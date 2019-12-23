package console

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/conf"
	"github.com/spf13/cobra"
	"path/filepath"
)

var logger = logs.GetLogger("console")

var client rpc.RpcClient

func onMessage(channel remoting.Channel, request *rpc.Request) *rpc.Response {
	if request.URL == "tail.logger" {
		fmt.Print(string(request.Body))
		return nil
	} else {
		return rpc.OK(channel, request)
	}
}

func preRune(cmd *cobra.Command, args []string) (err error) {
	sock := conf.Config.Server.Sock
	if sockServer, err := cmd.PersistentFlags().GetString("sock"); err != nil {
		return err
	} else if sockServer != "" {
		if absPath, err := filepath.Abs(sockServer); err != nil {
			return err
		} else {
			sock = "unix:/" + absPath
		}
	}
	logger.Debug("连接服务端sock: ", sock)
	client = rpc.NewClient(sock, onMessage, nil)
	if err = client.Start(); err != nil {
		logger.Warn("连接服务错误: ", err)
		return
	}
	return nil
}

func runPost(cmd *cobra.Command, args []string) {
	_ = client.Close()
}

var ConsoleCmd = &cobra.Command{
	Use: "console", Short: "管理端命令", Long: "管理端命令",
}

var Cmds = []*cobra.Command{ConsoleCmd,
	shutdownCmd,
	startCmd, statusCmd, stopCmd,
	listCmd, addCmd, deleteCmd, modifyCmd,
	detailCmd, tailCmd,
}

func init() {
	for _, cmd := range Cmds {
		cmd.PersistentFlags().StringP("sock", "s", "", "连接服务端sock地址.")
	}
	for i := 1; i < len(Cmds); i++ {
		ConsoleCmd.AddCommand(Cmds[i])
	}
}
