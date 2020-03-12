package daemon

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/kexec"
	"github.com/ihaiker/gokit/remoting"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Process struct {
	Pid    int      `json:"pid"`
	Status FSMState `json:"status"`

	*Program

	cmd *kexec.KCommand

	log *ProcessLogger

	retryLeft int

	statusListener FSMStatusListener

	stopC chan syscall.Signal

	Cpu float64 `json:"cpu"`
	Rss uint64  `json:"rss"`
}

func NewProcess(program *Program) *Process {
	return &Process{
		Status:    Ready,
		Program:   program,
		retryLeft: program.StartRetries,
	}
}

func (f *Process) initLogger() (err error) {
	if f.log, err = NewLogger(f.Program.Logger); err != nil {
		return
	}
	return
}
func (f *Process) GetStatus() FSMState {
	return f.Status
}

func (f *Process) GetLogger() *ProcessLogger {
	return f.log
}

func (f *Process) setState(newState FSMState) {
	if f.Status == newState {
		return
	}
	if newState == Running {
		f.Pid = f.cmd.Process.Pid
	} else {
		f.Pid = 0
	}

	if f.statusListener != nil {
		f.statusListener(FSMStatusEvent{Process: f, FromStatus: f.Status, ToStatus: newState})
	}
	f.Status = newState
}

func (p *Process) startedCallback(determinedResult chan *Process) {
	defer func() { _ = recover() }()
	if determinedResult != nil {
		determinedResult <- p
		close(determinedResult)
	}
}

func (p *Process) startCommand(determinedResult chan *Process) (err error) {
	if p.Status == Starting || p.Status == Running {
		return errors.New(fmt.Sprintf("the program(%s) is %s", p.Program.Name, p.Status))
	}
	startCmd := p.Program.Start

	logger.Debugf("start program: program=%s, cmd=%s, args=%s", p.Program.Name, startCmd.Command, strings.Join(startCmd.Args, " "))
	if p.cmd, err = p.buildCommand(startCmd); err != nil {
		return
	}

	p.setState(Starting)
	if err = p.cmd.Start(); err != nil {
		logger.Debug("start cmd(%s) error: %s", startCmd.Command, err)
		p.setState(Fail)
		return
	}

	p.stopC = make(chan syscall.Signal, 1)

	if p.Program.IsForeground() {
		startCheck := make(chan bool)
		go func() {
			exitErrChan := Async(p.cmd.Wait)
			select {
			case exit := <-exitErrChan:
				//关掉检查，这个时候确定已经出问题了
				safeCloseBool(startCheck)

				if p.Status == Stopping { //执行关闭动作直接退出
					return
				}

				if exit == nil {
					logger.Warnf("program(%s) is daemon? do not retry!", p.Program.Name)
				} else {
					logger.Warnf("program(%s) %s ", p.Program.Name, exit.Error())
				}

				if p.retryLeft > 0 {
					p.setState(RetryWait)
					go p.waitNextRetry(determinedResult)
					return
				} else {
					p.retryLeft = p.Program.StartRetries
					p.setState(Fail)
					logger.Warnf("program(%s) exit too quick, status -> fail", p.Program.Name)
					p.startedCallback(determinedResult)
					safeCloseSig(p.stopC)
				}
			case <-p.stopC:
				if p.Status == Starting || p.Status == Running {
					safeCloseBool(startCheck)
					_ = p.stopCommand()
				}
			}
		}()

		go func() {
			select {
			case <-startCheck:
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
							p.setState(Fail)
							if p.retryLeft > 0 {
								go p.waitNextRetry(determinedResult)
							} else {
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
						timer.Stop()
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
	if p.Status != Starting && p.Status != Running {
		return errors.New(fmt.Sprintf("the program(%s) is %s", p.Program.Name, p.Status))
	}

	p.setState(Stopping)
	safeCloseSig(p.stopC)

	p.retryLeft = p.Program.StartRetries

	if p.Program.IsForeground() {

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
			p.setState(Fail)
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

		p.setState(Stoped)
	}
	return nil
}

func (p *Process) waitNextRetry(determinedResult chan *Process) {
	if p.retryLeft <= 0 {
		p.retryLeft = p.Program.StartRetries
		p.setState(Fail)
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

	runCmd := command.Command
	if strings.HasPrefix(runCmd, "./") {
		abswd, _ := filepath.Abs(p.Program.WorkDir)
		runCmd = abswd + "/" + runCmd[2:]
	}

	cmd = kexec.Command(runCmd, args...)

	cmd.Env = os.Environ()
	if p.Program.Envs != nil && len(p.Program.Envs) > 0 {
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

	cmd.Stdout = p.log
	cmd.Stderr = p.log

	return
}

func (p *Process) refresh() {
	if p.Pid == 0 {
		return
	}
	var err error
	if p.Cpu, p.Rss, err = GetProcessInfo(int32(p.Pid)); err != nil {
		logger.Debug("get group info error: ", err)
	}
}

//释放所有资源
func (p *Process) Freed() {
	if p.log != nil {
		_ = p.log.Close()
	}
	p.Cpu = 0
	p.Rss = 0
}

func safeCloseSig(c chan syscall.Signal) {
	defer func() { recover() }()
	close(c)
}

func safeCloseBool(c chan bool) {
	defer func() { recover() }()
	close(c)
}
