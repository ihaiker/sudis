package server

import (
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"strings"
)

type masterClient interface {
	Start() error
	Stop()
	Notify(request *rpc.Request)
}

type emptyMasterClient struct{}

func (self *emptyMasterClient) Start() error {
	return nil
}

func (self *emptyMasterClient) Stop() {}

func (self *emptyMasterClient) Notify(req *rpc.Request) {

}

func newMasterClient(masterAddress, securityToken string, dm *daemon.DaemonManager) masterClient {
	if masterAddress == "" {
		return &emptyMasterClient{}
	} else if strings.HasPrefix(masterAddress, "tcp://") {
		return newTcpMasterClient(masterAddress[6:], securityToken, dm)
	} else {
		return newWSMasterClient(masterAddress, securityToken, dm)
	}
}
