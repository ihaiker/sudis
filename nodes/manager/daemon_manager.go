package manager

import (
	"encoding/json"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"strconv"
	"time"
)

type joinedDaemonManger struct {
	rpc.RpcServer
	key, address string
}

func NewManager(key, address string, server rpc.RpcServer) *cluster.NamedDomainManager {
	return &cluster.NamedDomainManager{
		Name: key, Address: address,
		Status: dao.NodeStatusOnline,
		Manager: &joinedDaemonManger{
			RpcServer: server, key: key, address: address,
		},
	}
}

func (self *joinedDaemonManger) Start() error {
	return nil
}

func (self *joinedDaemonManger) Stop() error {
	if c, has := self.GetChannel(self.address); has {
		return c.Close()
	}
	return nil
}

func makeRequest(cmd string, bodys ...string) *rpc.Request {
	req := &rpc.Request{URL: cmd}
	if len(bodys) > 0 {
		req.Body, _ = json.Marshal(bodys)
	}
	return req
}

func (self *joinedDaemonManger) sendRequest(req *rpc.Request) *rpc.Response {
	return self.Send(self.address, req, time.Second*7)
}

func (self *joinedDaemonManger) GetProcess(name string) (process *daemon.Process, err error) {
	req := makeRequest("detail", name)
	if resp := self.sendRequest(req); resp.Error != nil {
		err = resp.Error
	} else {
		process = new(daemon.Process)
		err = json.Unmarshal(resp.Body, process)
	}
	return
}

func (self *joinedDaemonManger) AddProgram(program *daemon.Program) error {
	req := makeRequest("add", program.JSON())
	resp := self.sendRequest(req)
	return resp.Error
}

func (self *joinedDaemonManger) RemoveProgram(name string, skip bool) error {
	req := makeRequest("delete", name)
	req.Header("skip", strconv.FormatBool(skip))
	resp := self.sendRequest(req)
	return resp.Error
}

func (self *joinedDaemonManger) ModifyProgram(program *daemon.Program) error {
	req := makeRequest("modify")
	req.Body = program.JSONByte()
	resp := self.sendRequest(req)
	return resp.Error
}

func (self *joinedDaemonManger) ListProgramNames() ([]string, error) {
	process, err := self.ListProcess()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(process))
	for i, p := range process {
		names[i] = p.Name
	}
	return names, err
}

func (self *joinedDaemonManger) ListProcess() (process []*daemon.Process, err error) {
	req := makeRequest("list")
	req.Header("inspect", "true")
	req.Header("all", "true")

	if resp := self.sendRequest(req); resp.Error != nil {
		err = resp.Error
	} else {
		process = make([]*daemon.Process, 0)
		err = json.Unmarshal(resp.Body, &process)
	}
	return
}

func (self *joinedDaemonManger) SetStatusListener(lis daemon.FSMStatusListener) {

}

func (self *joinedDaemonManger) StartProgram(name string, determinedResult chan *daemon.Process) error {
	req := makeRequest("start", name)
	self.Async(self.address, req, time.Second*30, func(response *rpc.Response) {
		errors.Try(func() {
			determinedResult <- self.MustGetProcess(name)
		})
	})
	return nil
}

func (self *joinedDaemonManger) StopProgram(name string) error {
	req := makeRequest("stop", name)
	resp := self.sendRequest(req)
	return resp.Error
}

func (self *joinedDaemonManger) MustGetProcess(name string) *daemon.Process {
	p, err := self.GetProcess(name)
	errors.Assert(err, "Get Process")
	return p
}

func (self *joinedDaemonManger) MustAddProgram(program *daemon.Program) {
	errors.Assert(self.AddProgram(program), "Add Program")
}

func (self *joinedDaemonManger) MustRemoveProgram(name string, skip bool) {
	errors.Assert(self.RemoveProgram(name, skip), "remove program")
}

func (self *joinedDaemonManger) MustModifyProgram(program *daemon.Program) {
	errors.Assert(self.ModifyProgram(program), "modify program")
}

func (self *joinedDaemonManger) MustListProgramNames() []string {
	names, err := self.ListProgramNames()
	errors.Assert(err, "list program names")
	return names
}

func (self *joinedDaemonManger) MustListProcess() []*daemon.Process {
	ps, err := self.ListProcess()
	errors.Assert(err, "list process")
	return ps
}

func (self *joinedDaemonManger) MustStartProgram(name string, determinedResult chan *daemon.Process) {
	err := self.StartProgram(name, determinedResult)
	errors.Assert(err, "start program ", name)
}

func (self *joinedDaemonManger) MustStopProgram(name string) {
	err := self.StopProgram(name)
	errors.Assert(err, "stop program ", name)
}

func (self *joinedDaemonManger) ModifyTag(name string, add bool, tag string) error {
	req := makeRequest("tag", name, tag)
	if !add {
		req.Header("delete", "true")
	}
	resp := self.sendRequest(req)
	return resp.Error
}

func (self *joinedDaemonManger) MustModifyTag(name string, add bool, tag string) {
	errors.Assert(self.ModifyTag(name, add, tag))
}

func (self *joinedDaemonManger) SubscribeLogger(uid string, name string, tail daemon.TailLogger, firstLine int) error {
	req := makeRequest("tail", name, "true", uid)
	req.Header("num", strconv.Itoa(firstLine))
	response := self.sendRequest(req)
	return response.Error
}

func (self *joinedDaemonManger) UnSubscribeLogger(uid, name string) error {
	resp := self.sendRequest(makeRequest("tail", name, "false", uid))
	return resp.Error
}
