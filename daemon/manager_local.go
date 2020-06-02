package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/logs"
	. "github.com/ihaiker/sudis/libs/errors"
	"os"
	"path/filepath"
	"sort"
)

var logger = logs.GetLogger("daemon")

type local struct {
	dataDir        *files.File
	group          ProcessGroup
	maxId          *atomic.AtomicUint64
	statusListener FSMStatusListener
	nodeKey        string
}

func NewDaemonManager(dirPath, nodeKey string) *local {
	dir := files.New(dirPath)
	return &local{
		dataDir: dir, nodeKey: nodeKey,
		group:          []*Process{},
		maxId:          atomic.NewAtomicUint64(1),
		statusListener: StdoutStatusListener,
	}
}

func (self *local) MaxId() uint64 {
	return self.maxId.IncrementAndGet(1)
}

func (self *local) SetStatusListener(lis FSMStatusListener) {
	self.statusListener = lis
	for _, process := range self.group {
		process.statusListener = self.statusListener
	}
}

func (self *local) notifyStatus(process *Process, oldStatus, newStatus FSMState) {
	if self.statusListener != nil {
		self.statusListener(FSMStatusEvent{
			Process:    process,
			FromStatus: oldStatus,
			ToStatus:   newStatus,
		})
	}
}

func (self *local) Start() error {
	if !self.dataDir.Exist() {
		if err := self.dataDir.Mkdir(); err != nil {
			return errors.Wrap(err, "can't use the dataDir: "+self.dataDir.GetPath())
		}
	} else if self.dataDir.IsFile() {
		return errors.Wrap(os.ErrNotExist, "can't use the dataDir: "+self.dataDir.GetPath())
	}

	logger.Info("start sudis daemon manager")

	programConfigs, _ := self.dataDir.List()
	for _, programConfig := range programConfigs {
		program := NewProgram()
		if bs, err := programConfig.ToBytes(); err != nil {
			logger.Warn("can't load config file: ", programConfig.GetPath())
		} else if err := json.Unmarshal(bs, program); err != nil {
			logger.Warn("Unmarshal program config error:", err)
		} else {
			if program.Id > self.maxId.Get() {
				self.maxId.Set(program.Id)
			}
			_ = self.AddProgram(program)
		}
	}

	sort.Sort(self.group)

	for _, pro := range self.group {
		if pro.Program.AutoStart {
			logger.Info("auto start program: ", pro.Program.Name)
			if err := pro.startCommand(nil); err != nil {
				logger.Warn("auto start program error:", err)
			}
		}
	}
	return nil
}

func (self *local) Stop() error {
	for _, process := range self.group {
		if process.GetStatus().IsRunning() {
			if err := process.stopCommand(); err != nil {
				logger.Debug("stop program", process.Program.Name, " error:", err)
			}
		}
		process.Freed()
	}
	logger.Info("stop sudis daemon manager")
	return nil
}

func (self *local) AddProgram(program *Program) error {
	program.Node = self.nodeKey
	if _, err := self.GetProcess(program.Name); err != ErrProgramNotFound {
		return ErrProgramExists
	}
	oldStatus := Ready
	if program.Id == 0 {
		program.Id = self.MaxId()
		oldStatus = ""
	}
	process := NewProcess(program)
	if err := process.initLogger(); err != nil {
		return err
	}
	process.statusListener = self.statusListener
	self.group = append(self.group, process)
	self.notifyStatus(process, oldStatus, process.Status)
	return self.WriteConfig(program.Name)
}

func (self *local) ListProcess() ([]*Process, error) {
	for _, d := range self.group {
		if d.Status.IsRunning() {
			d.refresh()
		}
	}
	return self.group, nil
}

func (self *local) GetProcess(name string) (p *Process, err error) {
	if p, err = self.group.Get(name); err == nil {
		if p.Status.IsRunning() {
			p.refresh()
		}
	}
	return
}

func (self *local) StartProgram(name string, determinedResult chan *Process) error {
	if p, err := self.GetProcess(name); err != nil {
		return err
	} else {
		return p.startCommand(determinedResult)
	}
}

func (self *local) StopProgram(name string) error {
	if p, err := self.GetProcess(name); err != nil {
		return err
	} else {
		return p.stopCommand()
	}
}

//skip 如果程序是运行状态不停止
func (self *local) RemoveProgram(name string, skip bool) error {
	if p, err := self.GetProcess(name); err != nil {
		return err
	} else {
		if p.GetStatus().IsRunning() {
			if !skip {
				return errors.New("程序正在运行")
			}
		}
		p.Freed()
		if err = self.group.Remove(name); err == nil {
			_ = os.Remove(filepath.Join(self.dataDir.GetPath(), fmt.Sprintf("%s.json", name)))
			self.notifyStatus(p, p.Status, "")
		}
		return err
	}
}

func (self *local) ModifyProgram(program *Program) error {
	if p, err := self.GetProcess(program.Name); err != nil {
		return err
	} else if p.GetStatus().IsRunning() {
		return errors.New("cant modify running program")
	} else {
		p.Freed()
		program.Id = p.Program.Id
		p.Program = program
		self.notifyStatus(p, Ready, p.Status)
		if err := p.initLogger(); err != nil {
			return err
		}
		return self.WriteConfig(program.Name)
	}
}

func (self *local) ModifyTag(name string, add bool, tag string) error {
	process, err := self.GetProcess(name)
	if err == nil {
		if add {
			process.Tags.Add(tag)
		} else {
			process.Tags.Remove(tag)
		}
	}
	return self.WriteConfig(name)
}

func (self *local) ListProgramNames() ([]string, error) {
	return self.group.Names(), nil
}

func (self *local) WriteConfig(name string) error {
	p, err := self.GetProcess(name)
	if err != nil {
		return err
	}

	program := p.Program
	file := files.New(filepath.Join(self.dataDir.GetPath(), program.Name+".json"))
	if file.Exist() {
		if err = file.Remove(); err != nil {
			return err
		}
	}
	if w, err := file.GetWriter(false); err != nil {
		return errors.New(fmt.Sprintf("\nwrite dm config %s error: %s", program.Name, err))
	} else {
		defer func() { _ = w.Close() }()
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		if err := encoder.Encode(program); err != nil {
			return errors.New(fmt.Sprintf("\nwrite dm config %s error: %s", program.Name, err))
		}
		return nil
	}
}

func (self *local) MustGetProcess(name string) *Process {
	p, err := self.GetProcess(name)
	errors.Assert(err, "Get Process")
	return p
}

func (self *local) MustAddProgram(program *Program) {
	errors.Assert(self.AddProgram(program), "Add Program")
}

func (self *local) MustRemoveProgram(name string, skip bool) {
	errors.Assert(self.RemoveProgram(name, skip), "remove program")
}

func (self *local) MustModifyProgram(program *Program) {
	errors.Assert(self.ModifyProgram(program), "modify program")
}

func (self *local) MustModifyTag(name string, add bool, tag string) {
	errors.Assert(self.ModifyTag(name, add, tag), "modify program tag")
}

func (self *local) MustListProgramNames() []string {
	names, err := self.ListProgramNames()
	errors.Assert(err, "list program names")
	return names
}

func (self *local) MustListProcess() []*Process {
	ps, err := self.ListProcess()
	errors.Assert(err, "list process")
	return ps
}

func (self *local) MustStartProgram(name string, determinedResult chan *Process) {
	err := self.StartProgram(name, determinedResult)
	errors.Assert(err, "start program ", name)
}

func (self *local) MustStopProgram(name string) {
	err := self.StopProgram(name)
	errors.Assert(err, "stop program ", name)
}

func (self *local) SubscribeLogger(uid string, name string, tail TailLogger, firstLine int) error {
	if process, err := self.GetProcess(name); err != nil {
		return err
	} else {
		process.GetLogger().Tail(uid, tail, firstLine)
		return nil
	}
}

func (self *local) UnSubscribeLogger(uid, name string) error {
	if process, err := self.GetProcess(name); err != nil {
		return err
	} else {
		process.GetLogger().CtrlC(uid)
		return nil
	}
}
