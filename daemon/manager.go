package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/files"
	"sort"
)

type DaemonManager struct {
	dir            *files.File
	process        ProcessList
	idx            *atomic.AtomicUint64
	statusListener FSMStatusListener
}

func (self *DaemonManager) MaxId() uint64 {
	return self.idx.IncrementAndGet(1)
}

func (self *DaemonManager) SetStatusListener(lis FSMStatusListener) {
	self.statusListener = lis
	for _, process := range self.process {
		process.statusListener = self.statusListener
	}
}

func (self *DaemonManager) Start() error {
	if !self.dir.IsDir() {
		if err := self.dir.Mkdir(); err != nil {
			return errors.New("can't use the dir: " + self.dir.GetPath())
		}
	}
	logger.Info("start sudis daemon manager")

	programConfigs, _ := self.dir.List()
	for _, programConfig := range programConfigs {
		program := NewProgram()
		if bs, err := programConfig.ToBytes(); err != nil {
			logger.Warn("can't load config file: ", programConfig.GetPath())
		} else if err := json.Unmarshal(bs, program); err != nil {
			logger.Warn("Unmarshal program config error:", err)
		} else {
			if program.Id > self.idx.Get() {
				self.idx.Set(program.Id)
			}
			_ = self.AddProgram(program)
		}
	}

	sort.Sort(self.process)

	for _, pro := range self.process {
		if pro.Program.AutoStart {
			logger.Debug("auto start program: ", pro.Program.Name)
			if err := pro.startCommand(nil); err != nil {
				logger.Warn("auto start program error:", err)
			}
		}
	}
	return nil
}

func (self *DaemonManager) Close() error {
	self.Stop()
	return nil
}

func (self *DaemonManager) Stop() {
	for _, process := range self.process {
		if process.GetStatus().IsRunning() {
			if err := process.stopCommand(); err != nil {
				logger.Debug("stop program", process.Program.Name, " error:", err)
			}
		}
		process.Freed()
	}
	logger.Info("stop sudis daemon manager")
}

func (self *DaemonManager) AddProgram(program *Program) error {
	if _, err := self.GetProgram(program.Name); err != ErrNotFound {
		return ErrExists
	}

	if program.Id == 0 {
		program.Id = self.MaxId()
	}
	process := NewProcess(program)
	if err := process.initLogger(); err != nil {
		return err
	}
	process.statusListener = self.statusListener
	self.process = append(self.process, process)
	if self.statusListener != nil {
		self.statusListener(process, "", process.State)
	}
	return nil
}

func (self *DaemonManager) ListProcess() []*Process {
	return self.process
}

func (self *DaemonManager) GetProgram(name string) (*Process, error) {
	return self.process.GetName(name)
}

func (self *DaemonManager) StartProgram(name string, determinedResult chan *Process) error {
	if p, err := self.GetProgram(name); err != nil {
		return err
	} else {
		return p.startCommand(determinedResult)
	}
}

func (self *DaemonManager) StopProgram(name string) error {
	if p, err := self.GetProgram(name); err != nil {
		return err
	} else {
		return p.stopCommand()
	}
}

//skip 如果程序是运行状态不停止
func (self *DaemonManager) RemoveProgram(name string, skip bool) error {
	if p, err := self.GetProgram(name); err != nil {
		return err
	} else {
		if p.GetStatus().IsRunning() {
			if !skip {
				return errors.New("程序正在运行")
			}
		}
		p.Freed()
		if err = self.process.Remove(name); err == nil {
			if self.statusListener != nil {
				self.statusListener(p, p.State, "")
			}
		}
		return err
	}
}

func (self *DaemonManager) ModifyProgram(program *Program) error {
	if p, err := self.GetProgram(program.Name); err != nil {
		return err
	} else if p.GetStatus().IsRunning() {
		return errors.New("cant modify running program")
	} else {
		p.Freed()
		program.Id = p.Program.Id
		p.Program = program
		return p.initLogger()
	}
}

func (self *DaemonManager) ListProgramNames() []string {
	return self.process.Names()
}

func (self *DaemonManager) WriteConfig(name string) error {
	p, err := self.GetProgram(name)
	if err != nil {
		return err
	}
	program := p.Program
	file := files.New(self.dir.GetPath() + "/" + program.Name + ".json")
	if file.Exist() {
		if err = file.Remove(); err != nil {
			return err
		}
	}
	if w, err := file.GetWriter(false); err != nil {
		return errors.New(fmt.Sprintf("\nwrite dm config %s error: %s", program.Name, err))
	} else {
		defer func() { _ = w.Close() }()
		bs, _ := json.MarshalIndent(program, "", "    ") //这里最要记录ID
		if _, err := w.Write(bs); err != nil {
			return errors.New(fmt.Sprintf("\nwrite dm config %s error: %s", program.Name, err))
		}
		return nil
	}
}

func NewDaemonManager(dirPath string) *DaemonManager {
	dir := files.New(dirPath)
	return &DaemonManager{
		dir: dir, process: []*Process{},
		idx: atomic.NewAtomicUint64(1),
	}
}
