package master

import (
	"errors"
	"github.com/ihaiker/gokit/logs"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/server"
	"github.com/ihaiker/sudis/master/server/eventbus"
	"github.com/ihaiker/sudis/master/server/http"
	"github.com/ihaiker/sudis/master/server/tcp"
	"time"
)

var logger = logs.GetLogger("master")

func StartAt(wait *runtimeKit.SignalListener) error {
	if conf.Config.Master.Band == "" && !conf.Config.Master.EnableWS {
		return errors.New("master.band or master enableWs mast config one")
	}

	//初始化应用状态
	if err := dao.CreateEngine(conf.Config.Master.Database); err == nil {
		dao.NodeDao.Ready()
		dao.ProgramDao.Ready()
	} else {
		return err
	}

	servers := []server.MasterServer{}
	api := server.NewApiWrapper()
	if conf.Config.Master.Band != "" {
		tcpServer := tcp.NewMasterTcpServer(conf.Config.Master.Band, func() {
			logger.Debug("程序接收关闭命令")
			wait.Stop()
		}, api)
		servers = append(servers, tcpServer)
	}
	httpServer := http.NewMasterHttpServer(conf.Config.Master.Http, conf.Config.Master.EnableWS, api)
	servers = append(servers, httpServer)

	eventbus.Service.AddListener(&eventbus.CommandListener{Api: api})
	servers = append(servers, eventbus.Service)

	for _, server := range servers {
		if err := server.Start(); err != nil {
			return err
		}
	}

	wait.OnClose(func() {
		for _, server := range servers {
			if err := server.Stop(); err != nil {
				logger.Warn("close server error:", err)
			}
		}
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
