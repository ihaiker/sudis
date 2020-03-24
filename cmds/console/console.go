package console

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/libs/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
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
	sockCfg := filepath.Join(config.Config.DataPath, "sudis.sock")
	if param := viper.GetString("sock"); param != "" {
		sockCfg = param
	}
	sock, _ := filepath.Abs(sockCfg)
	sock = "unix:/" + sock
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

func makeRequest(cmd string, body ...string) *rpc.Request {
	request := new(rpc.Request)
	request.URL = cmd
	if node := viper.GetString("node"); node != "" {
		request.Header("node", node)
	}
	if len(body) > 0 {
		request.Body, _ = json.Marshal(body)
	}
	return request
}

func sendRequest(request *rpc.Request, disablePrintln ...bool) *rpc.Response {
	seconds := viper.GetDuration("timeout")
	request.Header("timeout", fmt.Sprintf("%.0f", seconds.Seconds()))
	resp := client.Send(request, seconds)
	if len(disablePrintln) > 0 && disablePrintln[0] {
		//
	} else {
		if resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
	}
	return resp
}

var ConsoleCommands = &cobra.Command{
	Use: "console", Short: "管理端命令", Long: "管理端命令", Aliases: []string{"cli"},
}

var Commands = []*cobra.Command{
	startCmd, statusCmd, stopCmd,
	listCmd, addCmd, deleteCmd, modifyCmd,
	detailCmd, tailCmd, tagCommand, joinCmd,
}

func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("sock", "s", "", "连接服务端sock地址.(default: ${data-path}/sudis.sock)")
	cmd.PersistentFlags().DurationP("timeout", "t", time.Second*15, "wait timeout")
	cmd.PersistentFlags().StringP("node", "k", "", "执行的节点")
}

func init() {
	for _, cmd := range Commands {
		addFlags(cmd)
		_ = viper.BindPFlags(cmd.PersistentFlags())
	}
	addFlags(ConsoleCommands)
	ConsoleCommands.AddCommand(Commands...)
	_ = viper.BindPFlags(ConsoleCommands.PersistentFlags())
}
