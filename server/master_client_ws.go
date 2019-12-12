package server

import (
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
)

type wsMasterClient struct {
}

func newWSMasterClient(address, securityToken string, dm *daemon.DaemonManager) *wsMasterClient {
	return &wsMasterClient{}
}

func (self *wsMasterClient) Start() error {
	return nil
}

func (self *wsMasterClient) Stop() {

}

func (self *wsMasterClient) Notify(req *rpc.Request) {

}
