package node

import (
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting/rpc"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/ihaiker/sudis/daemon"
	. "github.com/ihaiker/sudis/libs/config"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/command"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/ihaiker/sudis/nodes/http"
	"github.com/ihaiker/sudis/nodes/join"
	"github.com/ihaiker/sudis/nodes/manager"
	"github.com/ihaiker/sudis/nodes/notify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func Start() (err error) {
	defer errors.Catch(func(re error) {
		err = re
	})
	logs.SetDebugMode(Config.Debug)
	logs.Info("Using config file:", viper.ConfigFileUsed())

	//启动数据库
	errors.Assert(dao.CreateEngine(Config.DataPath, Config.Database))

	signal := runtimeKit.NewListener()

	clusterManger := makeDaemonManager(signal)

	//管理节点启动
	if Config.Manager != "" {
		signal.Add(manager.NewJoinServer(Config.Manager, Config.Salt, clusterManger))
	}

	joinManager := join.New(Config.Key, Config.Salt)

	//open api and web ui
	if Config.Address != "" {
		signal.Add(http.NewHttpServer(Config.Address, Config.DisableWebUI, clusterManger, joinManager))
	}

	makeSockConsoleListener(signal, clusterManger, joinManager)
	makeJoinManager(signal, joinManager)

	signal.Add(notify.New(Config.NotifySynchronize, clusterManger, joinManager))

	return signal.WaitTimeout(Config.MaxWaitTimeout)
}

func makeDaemonManager(signal *runtimeKit.SignalListener) *cluster.DaemonManager {
	localDaemon := daemon.NewDaemonManager(filepath.Join(Config.DataPath, "programs"), Config.Key)
	clusterManger := cluster.NewDaemonManger(Config.Key, Config.Salt, localDaemon)
	signal.Add(clusterManger)
	return clusterManger
}

func makeSockConsoleListener(signal *runtimeKit.SignalListener, daemonManger *cluster.DaemonManager, joinManager *join.ToJoinManager) {
	sock, _ := filepath.Abs(filepath.Join(Config.DataPath, "sudis.sock"))
	_ = os.Remove(sock)
	sockAddress := fmt.Sprintf("unix:/%s", sock)
	joinManager.OnRpcMessage = command.MakeCommand(daemonManger, joinManager)
	signal.Add(rpc.NewServer(sockAddress, joinManager.OnRpcMessage, nil))
}

func makeJoinManager(signal *runtimeKit.SignalListener, joinManager *join.ToJoinManager) {
	signal.Add(joinManager)
	signal.AddStart(func() error {
		for _, joinAttr := range Config.Join {
			addressAndToken := strings.SplitN(joinAttr, ",", 2)
			if len(addressAndToken) == 1 {
				if Config.Salt == "" {
					logs.Warnf("ignore to join %s. flag `--slat` and token is empty", joinAttr)
					continue
				}
				addressAndToken = append(addressAndToken, Config.Salt)
			}
			joinManager.MustJoinIt(addressAndToken[0], addressAndToken[1])
		}
		return nil
	})
}
