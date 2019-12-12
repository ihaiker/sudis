package daemon

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/kexec"
	"github.com/ihaiker/gokit/remoting"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
	"syscall"
	"time"
)

type Process struct {
	State FSMState `json:"state"`

	Program *Program `json:"program"`

	cmd *kexec.KCommand

	stdout io.Writer
	stderr io.Writer

	retryLeft int

	statusListener FSMStatusListener

	stopC chan syscall.Signal
}

func NewProcess(program *Program) *Process {
	return &Process{
		State:     Ready,
		Program:   program,
		retryLeft: program.StartRetries,
	}
}

func (f *Process) GetStatus() FSMState {
	return f.State
}

func (f *Process) setState(newState FSMState) {
	if f.State == newState {
		return
	}
	if f.statusListener != nil {
		f.statusListener(f, f.State, newState)
	}
	f.State = newState
}

func (p *Process) startedCallback(determinedResult chan *Process) {
	defer func() { recover() }()
	if determinedResult != nil {
		determinedResult <- p
		close(determinedResult)
	}
}

func (p *Process) startCommand(determinedResult chan *Process) (err error) {
	if p.State == Starting || p.State == Running {
		return errors.New(fmt.Sprintf("the program(%s) is %s", p.Program.Name, p.State))
	}
	startCmd := p.Program.Start

	logger.Debugf("start program: program=%s, cmd=%s, args=%s", p.Program.Name, startCmd.Command, strings.Join(startCmd.Args, " "))
	if p.cmd, err = p.buildCommand(startCmd); err != nil {
		return
	}

	p.setState(Starting)
	if err = p.cmd.Start(); err != nil {
		logger.Debug("start cmd(%s) error: %s", startCmd.Command, err)
		p.setState(Fatal)
		return
	}

	p.stopC = make(chan syscall.Signal, 1)

	if p.Program.IsForeground() {
		startCheck := make(chan bool)
		go func() {
			exitErrChan := Async(p.cmd.Wait)
			select {
			case exit := <-exitErrChan:
				p.setState(Fatal)
				if p.State == Starting {
					safeCloseBool(startCheck)
				}
				if exit == nil {
					logger.Warnf("program(%s) is daemon? do not retry!", p.Program.Name)
				} else {
					logger.Warnf("program(%s) %s ", p.Program.Name, exit.Error())
				}
				if p.retryLeft > 0 {
					p.waitNextRetry(determinedResult)
				} else {
					logger.Warnf("program(%s) exit too quick, status -> fatal", p.Program.Name)
					p.startedCallback(determinedResult)
					safeCloseSig(p.stopC)
				}
			case <-p.stopC:
				if p.State == Starting {
					safeCloseBool(startCheck)
				}
				_ = p.stopCommand()
			}
		}()

		go func() {
			select {
			case <-startCheck:
				p.startedCallback(determinedResult)
				return
			case <-time.After(time.Second * time.Duration(p.Program.StartDuration)):
				p.setState(Running)
				p.startedCallback(determinedResult)
				p.retryLeft = p.Program.StartRetries
				safeCloseBool(startCheck)
			}
		}()

	} else {
		err := p.cmd.Wait()
		//如果推出状态不正切且不忽略已经启动，说明出问题了
		if err != nil && !p.Program.IgnoreAlreadyStarted {
			p.waitNextRetry(determinedResult)
		} else {
			//如果不是上面的情况，则说明需要监控检查
			go func() {
				timer := time.NewTimer(time.Millisecond * 20)
				for {
					select {
					case <-timer.C:
						timer.Reset(time.Second * time.Duration(startCmd.CheckHealth.CheckTtl))
						if err := p.healthCheck(startCmd); err != nil {
							logger.Warnf("the program health is error: program=%s, num=%d, error=%s", p.Program.Name, p.Program.StartRetries-p.retryLeft, err)
							if p.retryLeft > 0 {
								p.waitNextRetry(determinedResult)
							} else {
								p.setState(Fatal)
								p.startedCallback(determinedResult)
								safeCloseSig(p.stopC)
							}
							return
						} else {
							p.retryLeft = p.Program.StartRetries
							p.setState(Running)
							p.startedCallback(determinedResult)
						}
					case <-p.stopC:
						_ = p.stopCommand()
						return
					}
				}
			}()
		}
	}
	return
}

func (p *Process) healthCheck(cmd *Command) error {
	logger.Debugf("health check program: name=%s, type=%s url=%s", p.Program.Name, cmd.CheckHealth.CheckMode, cmd.CheckHealth.CheckAddress)
	address := cmd.CheckHealth.CheckAddress
	if cmd.CheckHealth.CheckMode == HTTP || cmd.CheckHealth.CheckMode == HTTPS {
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			address = fmt.Sprintf("%s://%s", cmd.CheckHealth.CheckMode, address)
		}
		if cmd.CheckHealth.SecretToken != "" {
			if strings.Contains(address, "?") {
				address = address + "&" + cmd.CheckHealth.SecretToken
			} else {
				address = address + "?" + cmd.CheckHealth.SecretToken
			}
		}
		if resp, err := http.Get(address); err != nil {
			return err
		} else if resp.StatusCode/100 != 2 {
			return errors.New(resp.Status)
		}
	} else {
		if con, err := remoting.Dial(address); err != nil {
			return err
		} else {
			_ = con.Close()
		}
	}

	return nil
}

func (p *Process) stopCommand() error {
	if p.State != Starting && p.State != Running {
		return errors.New(fmt.Sprintf("the program(%s) is %s", p.Program.Name, p.State))
	}

	p.setState(Stopping)
	safeCloseSig(p.stopC)

	if p.Program.IsForeground() {
		p.retryLeft = p.Program.StartRetries

		//正常退出模式
		_ = p.cmd.Terminate(p.Program.StopSign)

		select {
		case <-Async(p.cmd.Wait):
			logger.Debugf("program(%s) quit normally", p.Program.Name)
		case <-time.After(time.Duration(p.Program.StopTimeout) * time.Second):
			logger.Debugf("program(%s) terminate all", p.Program.Name)
			_ = p.cmd.Terminate(syscall.SIGKILL)
		}

		_ = p.cmd.Wait() // This is OK, because Signal KILL will definitely work

		p.setState(Stoped)
	} else {
		stopCmd, err := p.buildCommand(p.Program.Stop)
		if err != nil {
			p.setState(Fatal)
			return errors.New(fmt.Sprintf("make program(%s) stop command error: %v", p.Program.Name, err))
		}

		if err = stopCmd.Start(); err != nil {
			return errors.New(fmt.Sprintf("exec program(%s) stop command error: %v", p.Program.Name, err))
		}

		//执行等待防止程序无法停止卡死
		select {
		case <-Async(stopCmd.Wait):
		case <-time.After(time.Duration(p.Program.StopTimeout) * time.Second):
			_ = stopCmd.Terminate(syscall.SIGKILL)
		}

		_ = stopCmd.Wait()
	}
	return nil
}

func (p *Process) waitNextRetry(determinedResult chan *Process) {
	if p.retryLeft <= 0 {
		p.setState(Fatal)
		return
	}
	p.retryLeft -= 1

	p.setState(RetryWait)
	logger.Debugf("retry program(%s) : %v", p.Program.Name, p.Program.StartRetries-p.retryLeft)

	select {
	case <-p.stopC:
		_ = p.stopCommand()
		p.startedCallback(determinedResult)

	case <-time.After(2 * time.Second):
		safeCloseSig(p.stopC)
		_ = p.startCommand(determinedResult) //can ignore,如果可以到这里说明之前已经验证过肯定不会返回error
	}
}

//创建启动命令
func (p *Process) buildCommand(command *Command) (cmd *kexec.KCommand, err error) {
	args := command.Args
	for idx, arg := range args {
		args[idx] = os.ExpandEnv(arg)
	}

	cmd = kexec.Command(command.Command, args...)

	cmd.Env = os.Environ()
	if p.Program.Envs != nil || len(p.Program.Envs) > 0 {
		programEvns := p.Program.Envs
		for idx, env := range programEvns {
			programEvns[idx] = os.ExpandEnv(env)
		}
		cmd.Env = append(cmd.Env, programEvns...)
	}

	dir := os.TempDir()
	if p.Program.WorkDir != "" {
		dir = p.Program.WorkDir
	}
	cmd.Dir = os.ExpandEnv(dir)

	currentUser, _ := user.Current()
	if p.Program.User != "" {
		if currentUser, err = user.Lookup(currentUser.Username); err != nil {
			return
		}
	}

	cmd.Env = append(cmd.Env, "HOME="+currentUser.HomeDir, "USER="+currentUser.Username)

	//FIXME 翻入输出缓存
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	return
}

func safeCloseSig(c chan syscall.Signal) {
	defer func() { recover() }()
	close(c)
}

func safeCloseBool(c chan bool) {
	defer func() { recover() }()
	close(c)
}
