package cluster

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/daemon"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/dao"
	"net"
	"time"
)

var logger = logs.GetLogger("cluster")

type (
	DaemonManager struct {
		key, salt      string
		domainManagers []*NamedDomainManager
		tails          map[string]daemon.TailLogger
		eventListener  daemon.FSMStatusListener
		nodeListener   NodeListener
	}
	//节点事件
	NodeEvent struct {
		Node   string         `json:"node"`
		Status dao.NodeStatus `json:"status"`
	}
	NodeListener func(event NodeEvent)

	NamedDomainManager struct {
		Name    string
		Address string
		Status  dao.NodeStatus
		daemon.Manager
	}
)

func NewDaemonManger(key, salt string, daemonManger daemon.Manager) *DaemonManager {
	return &DaemonManager{
		key: key, salt: salt,
		tails: map[string]daemon.TailLogger{},
		domainManagers: []*NamedDomainManager{
			{
				Name: key, Address: "127.0.0.1",
				Status: dao.NodeStatusOnline, Manager: daemonManger,
			},
		},
	}
}

func (self DaemonManager) Get(key string) (*NamedDomainManager, error) {
	for _, manger := range self.domainManagers {
		if key == manger.Name {
			return manger, nil
		}
	}
	return nil, ErrNodeNotFound
}

func (self DaemonManager) MustGet(key string) *NamedDomainManager {
	for _, manger := range self.domainManagers {
		if key == manger.Name {
			return manger
		}
	}
	panic(ErrNodeNotFound)
}

func (self *DaemonManager) Start() error {
	for _, manager := range self.domainManagers {
		if err := manager.Start(); err != nil {
			return err
		}
		manager.SetStatusListener(self.OnStatusEvent)
	}
	return nil
}

func (self *DaemonManager) Stop() error {
	for _, manager := range self.domainManagers {
		if err := manager.Stop(); err != nil {
			logger.Warnf("stop %s daemon error", manager.Name)
		}
	}
	return nil
}

func (self *DaemonManager) GetProgramNum(key string) int {
	defer errors.Catch()
	return len(self.MustGet(key).MustListProcess())
}

func (self *DaemonManager) GetStatus(key string) (status dao.NodeStatus) {
	defer errors.Catch()
	status = dao.NodeStatusOutline
	m := self.MustGet(key)
	status = m.Status
	return
}

func match(value, search string) bool {
	return search == "" || value == search
}

func matchTag(tags []string, tag string) bool {
	if tag == "" {
		return true
	}
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func (self *DaemonManager) ListPrograms(name string, node string, tag string, status string, page int, limit int) *dao.Page {
	data := make([]*daemon.Process, 0)

	for _, manger := range self.domainManagers {
		if manger.Status == dao.NodeStatusOutline {
			continue
		}
		if !match(manger.Name, node) {
			continue
		}
		process, err := manger.ListProcess()
		if err != nil {
			logger.Warnf("list %s process error: %s", manger.Name, err)
			continue
		}

		for _, process := range process {
			if !match(process.Program.Name, name) {
				continue
			}
			if !matchTag(process.Program.Tags, tag) {
				continue
			}
			if !match(process.Status.String(), status) {
				continue
			}
			data = append(data, process)
		}
	}

	start := min((page-1)*limit, len(data))
	end := min(start+limit, len(data))

	return &dao.Page{
		Total: int64(len(data)), Data: data[start:end],
		Page: page, Limit: limit,
	}
}

func (self *DaemonManager) ModifyProgramTag(name string, node string, tag string, add bool) (err error) {
	defer errors.Catch(func(re error) {
		err = re
	})
	d := self.MustGet(node)
	err = d.ModifyTag(name, add, tag)
	return
}

func (self *DaemonManager) Command(node, name, command string, timeout time.Duration) (err error) {
	defer errors.Catch(func(re error) {
		err = re
	})

	d := self.MustGet(node)

	switch command {
	case "delete":
		d.MustRemoveProgram(name, false)
	case "start":
		{
			result := make(chan *daemon.Process, 1)
			defer errors.Try(func() { close(result) })
			d.MustStartProgram(name, result)

			select {
			case out, has := <-result:
				errors.True(has && out.Status == daemon.Running, "程序启动失败:", out.Status)
			case <-time.After(timeout):
			}
		}
	case "stop":
		d.MustStopProgram(name)
	case "restart":
		errors.Assert(self.Command(node, name, "stop", timeout))
		errors.Assert(self.Command(node, name, "start", timeout))
	}
	return
}

func (self *DaemonManager) MustCommand(node, name, command string) {
	errors.Assert(self.Command(node, name, command, time.Second*7))
}

func (self *DaemonManager) QueryNode() []*dao.Node {
	nodes, err := dao.NodeDao.List()
	errors.Assert(err)

	for _, node := range nodes {
		node.ProgramNum = self.GetProgramNum(node.Key)
		node.Status = self.GetStatus(node.Key)
	}
	nodes = append(nodes, &dao.Node{
		Tag: "本地", Key: self.key, Ip: "127.0.0.1", Address: "本地",
		ProgramNum: self.GetProgramNum(self.key),
		Status:     dao.NodeStatusOnline, Time: time.Now().Format("2006-01-02 15:04:05"),
	})
	return nodes
}

func (self DaemonManager) ModifyNodeToken(key string, token string) error {
	if key == self.key {
		return nil
	}
	return dao.NodeDao.ModifyToken(key, token)
}

func (self *DaemonManager) ModifyNodeTag(key string, tag string) error {
	return dao.NodeDao.ModifyTag(key, tag)
}

func (self *DaemonManager) GetProgram(node string, name string) (*daemon.Program, error) {
	d, err := self.Get(node)
	if err != nil {
		return nil, errors.Wrap(err, "节点不存在")
	}
	process, err := d.GetProcess(name)
	if err != nil {
		return nil, errors.Wrap(err, "获取托管程序异常")
	}
	return process.Program, nil
}

func (self *DaemonManager) AddProgram(node string, program *daemon.Program) error {
	if d, err := self.Get(node); err != nil {
		return err
	} else {
		return d.AddProgram(program)
	}
}

func (self *DaemonManager) MustAddProgram(node string, program *daemon.Program) {
	d := self.MustGet(node)
	d.MustAddProgram(program)
}

func (self *DaemonManager) MustModifyProgram(node string, name string, program *daemon.Program) {
	d := self.MustGet(node)
	program.Name = name
	program.Node = node
	d.MustModifyProgram(program)
}

func (self *DaemonManager) MustGetProcess(node string, name string) *daemon.Process {
	d := self.MustGet(node)
	return d.MustGetProcess(name)
}

var none = make([]*daemon.Process, 0)

func (self *DaemonManager) CacheAll() map[*dao.Node][]*daemon.Process {
	nodes := self.QueryNode()
	mapNodes := map[*dao.Node][]*daemon.Process{}
	for _, node := range nodes {
		d, err := self.Get(node.Key)
		if err != nil || d.Status == dao.NodeStatusOutline {
			mapNodes[node] = none
		} else {
			if mapNodes[node], err = d.ListProcess(); err != nil {
				logger.Warnf("获取节点 %s 程序列表异常：%s", node.Key, node.Tag)
			}
		}
	}
	return mapNodes
}

func (self *DaemonManager) SubscribeLogger(uid, node, name string, tail daemon.TailLogger, line int) (err error) {
	defer errors.Catch(func(re error) {
		err = re
	})
	d := self.MustGet(node)
	err = d.SubscribeLogger(uid, name, self.OnLogger, line)
	if err == nil {
		self.tails[uid] = tail
	}
	return
}

func (self *DaemonManager) UnsubscribeLogger(uid, node, name string) (err error) {
	defer errors.Catch(func(re error) { err = re })
	defer errors.Try(func() { delete(self.tails, uid) })

	d := self.MustGet(node)
	err = d.UnSubscribeLogger(uid, name)
	return
}

func (self *DaemonManager) OnLogger(uid, line string) {
	defer errors.Catch(func(re error) {
		logger.Warn(errors.Stack())
	})
	if tail, has := self.tails[uid]; has {
		tail(uid, line)
	}
}

func (self *DaemonManager) OnStatusEvent(event daemon.FSMStatusEvent) {
	defer errors.Catch(func(re error) {
		logger.Warn(errors.Stack())
	})
	if self.eventListener != nil {
		self.eventListener(event)
	}
}

func (self *DaemonManager) NodeJoin(key, address, timestamp, code string, dm *NamedDomainManager) error {
	token := self.salt
	if node, has, err := dao.NodeDao.Get(key); err == nil && has {
		token = node.Token
	}
	checkCode := fmt.Sprintf("%x", md5.Sum([]byte(timestamp+token+key)))
	if code != checkCode {
		return ErrToken
	}

	ip, _, _ := net.SplitHostPort(address)
	if err := dao.NodeDao.Join(ip, key); err != nil {
		logger.Warn("节点 %s,%s 加入异常：%s", key, address, err)
		return ErrToken
	}

	self.domainManagers = append(self.domainManagers, dm)
	dm.SetStatusListener(self.OnStatusEvent)

	if self.nodeListener != nil {
		self.nodeListener(NodeEvent{
			Node: key, Status: dao.NodeStatusOnline,
		})
	}
	return nil
}

func (self *DaemonManager) NodeLeave(key, address string) {
	ip, _, _ := net.SplitHostPort(address)
	if err := dao.NodeDao.Lost(ip, key); err != nil {
		logger.Warn("节点 %s,%s 离开异常：%s", key, ip, err)
	}
	idx := -1
	for i, manger := range self.domainManagers {
		if manger.Name == key {
			idx = i
		}
	}
	if idx != -1 {
		if self.nodeListener != nil {
			self.nodeListener(NodeEvent{
				Node:   self.domainManagers[idx].Name,
				Status: dao.NodeStatusOutline,
			})
		}
		self.domainManagers = append(self.domainManagers[0:idx], self.domainManagers[idx+1:]...)
	}
}

func (self DaemonManager) RemoveJoin(key string) error {
	if d, err := self.Get(key); err == nil {
		if err = d.Stop(); err != nil {
			logger.Warn("remove join error: ", err)
		}
	}
	return dao.NodeDao.Remove(key)
}

func (self *DaemonManager) SetStatusListener(statusListener daemon.FSMStatusListener) {
	self.eventListener = statusListener
}

func (self *DaemonManager) SetNodeListener(listener NodeListener) {
	self.nodeListener = listener
}
