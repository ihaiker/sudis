package eventbus

import (
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/server"
)

type EventListener interface {
	OnEvent(event *Event)
}

type CommandListener struct {
	Api server.Api
}

func (self *CommandListener) lostNode(node, key string) {
	logger.Info("节点丢失：node=", node, ",key=", key)
	_ = dao.NodeDao.Lost(node, key)
	dao.ProgramDao.Lost(key)
}

func (self *CommandListener) newNode(ip, key string) {
	logger.Info("新节点：ip=", ip, "，key=", key)
	_ = dao.NodeDao.Add(ip, key)

	if names, err := self.Api.ListProgramNames(key); err != nil {
		logger.Warn("list process error:", err)
		return
	} else {
		_ = dao.NodeDao.UpdateNodesProcessNumber(key, len(names))
		for _, name := range names {
			if p, err := self.Api.Get(key, name); err == nil {
				status := p.GetStatus()
				self.addOrModifyProgram(key, name, status)
			}
		}
	}
}

func (self *CommandListener) addOrModifyProgram(node, name string, status daemon.FSMState) {
	_, has, err := dao.ProgramDao.Get(name, node)
	if !has || err != nil {
		pro := &dao.Program{
			Name: name, Node: node, Status: status,
			Tags: []string{}, Time: dao.Timestamp(), Sort: 0,
		}
		if err = dao.ProgramDao.Add(pro); err != nil {
			logger.Warn("添加程序异常：node=", node, ",name=", name)
		}
	} else if err = dao.ProgramDao.UpdateStatus(node, name, status); err != nil {
		logger.Error("更新进程状态异常：node=", node, ",name=", name, ",status=", status, ",error:", err)
	} else {
		logger.Debug("更新进程状态：node=", node, ",name=", name, ",status=", status)
	}
}

func (self *CommandListener) updateNodeProgramNum(node string) {
	if names, err := self.Api.ListProgramNames(node); err != nil {
		logger.Warn("list process error:", err)
		return
	} else {
		_ = dao.NodeDao.UpdateNodesProcessNumber(node, len(names))
	}
}

func (self *CommandListener) programStatus(event *ProgramStatusEvent) {
	//添加或者删除
	if event.NewStatus == "" || event.OldStatus == "" {
		self.updateNodeProgramNum(event.Key)
	}
	if event.NewStatus == "" { //删除
		if err := dao.ProgramDao.Remove(event.Key, event.Name); err != nil {
			logger.Warn("删除程序异常：node=", event.Key, ",name=", event.Name, ",error=", err)
		}
	} else {
		self.addOrModifyProgram(event.Key, event.Name, daemon.FSMState(event.NewStatus))
	}
}

func (self *CommandListener) OnEvent(event *Event) {
	switch event.Name {
	case "NewNode":
		args := event.Value.([]string)
		self.newNode(args[0], args[1])
	case "LostNode":
		args := event.Value.([]string)
		self.lostNode(args[0], args[1])
	case "ProgramStatus":
		args := event.Value.(*ProgramStatusEvent)
		self.programStatus(args)
	}
}
