package server

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/daemon"
	"time"
)

type tcpMasterClient struct {
	client        rpc.RpcClient
	securityToken string
	dm            *daemon.DaemonManager
	started       bool
}

func newTcpMasterClient(address, securityToken string, dm *daemon.DaemonManager) *tcpMasterClient {
	client := &tcpMasterClient{
		started:       true,
		securityToken: securityToken,
		dm:            dm,
	}
	client.client = rpc.NewClient(address, MakeServerCommand(dm), func(channel remoting.Channel) {
		if client.started {
			logger.Info("与主控节点断开，重新连接")
			commons.SafeExec(client.Stop)
			for {
				if err := client.Start(); err == nil {
					return
				}
				time.Sleep(time.Second)
			}
		}
	})
	return client
}

func (self *tcpMasterClient) Notify(req *rpc.Request) {
	self.client.Oneway(req, time.Second*5)
}

func (self *tcpMasterClient) authRequest() *rpc.Request {
	req := new(rpc.Request)
	req.URL = "auth"
	timestamp := time.Now().Format("20060102150405")
	req.Header("timestamp", timestamp)
	req.Header("key", conf.Config.Server.Key)
	req.Body = []byte(fmt.Sprintf("%x", md5.Sum([]byte(timestamp+self.securityToken))))
	return req
}

func (self *tcpMasterClient) Start() (err error) {
	if err = self.client.Start(); err != nil {
		return
	}
	if resp := self.client.Send(self.authRequest(), time.Second*3); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (self *tcpMasterClient) Stop() {
	logger.Info("server master client close")
	self.started = false
	self.client.Shutdown()
	logger.Info("server master client closed")
}
