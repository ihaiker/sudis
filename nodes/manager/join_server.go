package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/nodes/cluster"
)

var logger = logs.GetLogger("manager")

type joinServer struct {
	salt string
	rpc.RpcServer
	dm *cluster.DaemonManager
}

func NewJoinServer(address, salt string, dm *cluster.DaemonManager) *joinServer {
	server := &joinServer{salt: salt, dm: dm}
	server.RpcServer = rpc.NewServer(
		address, server.onMessage, server.onClientClosed,
	)
	return server
}

func (self *joinServer) checkClientAuth(channel remoting.Channel, request *rpc.Request) *rpc.Response {
	if request.URL == "auth" {
		address := channel.GetRemoteAddress()
		timestamp, exits := request.GetHeader("timestamp")
		key, has := request.GetHeader("key")
		if exits && has {
			code := string(request.Body)
			daemonManger := NewManager(key, address, self)
			if err := self.dm.NodeJoin(key, address, timestamp, code, daemonManger); err == nil {
				channel.SetAttr("key", key)
				logger.Infof("node join key: %s, address: %s", key, address)
				return rpc.OK(channel, request)
			} else {
				logger.Warnf("node join key: %s, error: %s", key, err)
				return rpc.NewErrorResponse(request.ID(), err)
			}
		}
		return rpc.NewErrorResponse(request.ID(), errors.New("NoAuthHeader"))
	}

	if _, has := channel.GetAttr("key"); has {
		return nil
	} else {
		return rpc.NewErrorResponse(request.ID(), errors.New("ErrUnauthorized"))
	}
}

func (self *joinServer) onClientClosed(channel remoting.Channel) {
	if key, has := channel.GetAttr("key"); has {
		logger.Infof("node %s %s leave", key, channel.GetRemoteAddress())
		self.dm.NodeLeave(fmt.Sprint(key), channel.GetRemoteAddress())
	}
}

func (self *joinServer) onMessage(channel remoting.Channel, request *rpc.Request) *rpc.Response {
	if resp := self.checkClientAuth(channel, request); resp != nil {
		return resp
	}
	switch request.URL {
	case "program.status":
		event := new(daemon.FSMStatusEvent)
		if err := json.Unmarshal(request.Body, event); err != nil {
			logger.Warn("decoder event error ", err)
		} else {
			self.dm.OnStatusEvent(*event)
		}
		return rpc.OK(channel, request)
	case "tail.logger":
		id, _ := request.GetHeader("id")
		line := string(request.Body)
		self.dm.OnLogger(id, line)
	}
	return rpc.NewErrorResponse(request.ID(), rpc.ErrNotFount)
}
