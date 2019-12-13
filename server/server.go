package server

import (
	"encoding/json"
	"errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/daemon"
	"time"
)

var logger = logs.GetLogger("server")

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

	if err := dm.Start(); err != nil {
		return err
	}
	if err := server.Start(); err != nil {
		dm.Stop()
		return err
	}

	if err := client.Start(); err != nil {
		server.Shutdown()
		dm.Stop()
		return errors.New("connect master error: " + err.Error())
	}

	dm.SetStatusListener(func(process *daemon.Process, oldStatus, newStatus daemon.FSMState) {
		req := new(rpc.Request)
		req.URL = "ProgramStatus"
		req.Body, _ = json.Marshal([]string{process.Program.Name, string(oldStatus), string(newStatus)})
		client.Notify(req)
	})

	listener.PrependOnClose(func() {
		logger.Debug("关闭服务")
		client.Stop()
		server.Shutdown()
		dm.Stop()
	})

	return nil
}

func Start() error {
	listener := runtimeKit.NewListener()
	if err := StartAt(listener); err != nil {
		return err
	}
	return listener.WaitTimeout(time.Second * 10)
}
