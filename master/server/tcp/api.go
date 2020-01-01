package tcp

import (
	"encoding/json"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"strconv"
	"time"
)

type NodeApi struct {
	server *masterTcpServer
}

func (self *NodeApi) Online(node string) bool {
	_, has := self.server.GetChannel(node)
	return has
}

func (self *NodeApi) ListProgramNames(node string) ([]string, error) {
	req := new(rpc.Request)
	req.URL = "list"
	req.Body, _ = json.Marshal(&[]string{})
	programNames := make([]string, 0)
	if resp := self.server.Send(node, req, time.Second*7); resp.Error != nil {
		return nil, resp.Error
	} else if err := json.Unmarshal(resp.Body, &programNames); err != nil {
		return nil, err
	}
	return programNames, nil
}
func (self *NodeApi) ListProgram(node string) ([]*daemon.Process, error) {
	req := new(rpc.Request)
	req.URL = "list"
	req.Body, _ = json.Marshal(&[]string{"inspect"})
	processes := make([]*daemon.Process, 0)
	if resp := self.server.Send(node, req, time.Second*7); resp.Error != nil {
		return nil, resp.Error
	} else if err := json.Unmarshal(resp.Body, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}

func (self *NodeApi) Start(node, name string) error {
	request := new(rpc.Request)
	request.URL = "start"
	request.Body, _ = json.Marshal([]string{name})
	if resp := self.server.Send(node, request, time.Second*8); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *NodeApi) Stop(node, name string) error {
	request := new(rpc.Request)
	request.URL = "stop"
	request.Body, _ = json.Marshal([]string{name})
	if resp := self.server.Send(node, request, time.Second*7); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *NodeApi) Add(node string, program *daemon.Program) error {
	request := new(rpc.Request)
	request.URL = "add"
	programs := []string{program.JSON()}
	request.Body, _ = json.Marshal(programs)
	if resp := self.server.Send(node, request, time.Second*5); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *NodeApi) Remove(node, name string) error {
	request := new(rpc.Request)
	request.URL = "delete"
	request.Body, _ = json.Marshal([]string{name})
	request.Header("skip", "false")
	if resp := self.server.Send(node, request, time.Minute*7); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *NodeApi) Get(node, name string) (*daemon.Process, error) {
	request := new(rpc.Request)
	request.URL = "detail"
	request.Body, _ = json.Marshal([]string{name})

	if resp := self.server.Send(node, request, time.Second*7); resp.Error != nil {
		return nil, resp.Error
	} else {
		process := new(daemon.Process)
		if err := json.Unmarshal(resp.Body, process); err != nil {
			return nil, err
		} else {
			return process, nil
		}
	}
}

func (self *NodeApi) Modify(node, name string, program *daemon.Program) error {
	request := new(rpc.Request)
	request.URL = "modify"
	request.Body, _ = json.Marshal(program)
	if resp := self.server.Send(node, request, time.Second*7); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *NodeApi) TailLogger(node, name, id string, num int, consumer daemon.TailLogger) error {
	self.server.tails[id] = consumer
	request := new(rpc.Request)
	request.URL = "tail"
	request.Header("num", strconv.Itoa(num))
	request.Body, _ = json.Marshal([]string{name, "true", id})
	response := self.server.Send(node, request, time.Second*5)
	if response.Error != nil {
		delete(self.server.tails, id)
	}
	return response.Error
}

func (self *NodeApi) UnTailLogger(node, name, id string) error {
	request := new(rpc.Request)
	request.URL = "tail"
	request.Body, _ = json.Marshal([]string{name, "false", id})
	response := self.server.Send(node, request, time.Second*5)
	delete(self.server.tails, id)
	return response.Error
}
