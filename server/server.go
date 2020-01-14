package server

import (
	"encoding/json"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/daemon"
	"time"
)

var logger = logs.GetLogger("server")

type Services interface {
	Start() error
	Stop() error
}

func StartAt(listener *runtimeKit.SignalListener) error {
	dm := daemon.NewDaemonManager(conf.Config.Server.Dir)
	client := newMasterClient(conf.Config.Server.Master, conf.Config.Server.SecurityToken, dm)

	handler := MakeServerCommand(dm)
	server := rpc.NewServer(conf.Config.Server.Sock, func(channel remoting.Channel, request *rpc.Request) (resp *rpc.Response) {
		resp = rpc.NewResponse(request.ID())
		if request.URL == "shutdown" {
			resp.Body = []byte("OK")
			logger.Info("console shutdown")
			listener.Stop()
		} else {
			resp = handler(channel, request)
		}
		return resp
	}, nil)

	services := []Services{dm, client, server}

	for _, service := range services {
		if err := service.Start(); err != nil {
			return err
		}
	}

	dm.SetStatusListener(func(process *daemon.Process, oldStatus, newStatus daemon.FSMState) {
		req := new(rpc.Request)
		req.URL = "ProgramStatus"
		req.Body, _ = json.Marshal([]string{process.Program.Name, string(oldStatus), string(newStatus)})
		client.Notify(req)
	})

	listener.PrependOnClose(func() {
		logger.Debug("关闭服务")
		for i := len(services) - 1; i >= 0; i-- {
			_ = services[i].Stop()
		}
	})

	return nil
}

func Start() error {
	listener := runtimeKit.NewListener()
	if err := StartAt(listener); err != nil {
		return err
	}
	return listener.WaitTimeout(time.Second * 7)
}
