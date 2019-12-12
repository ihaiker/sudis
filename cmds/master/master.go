package master

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/master"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/spf13/cobra"
	"time"
)

var Cmd = &cobra.Command{
	Use: "master", Short: "管理控制端", Long: "管理控制端，提供报警服务，",
	RunE: func(cmd *cobra.Command, args []string) error {
		return master.Start()
	},
}

func authRequest() *rpc.Request {
	req := new(rpc.Request)
	req.URL = "auth"
	timestamp := time.Now().Format("20060102150405")
	req.Header("timestamp", timestamp)
	req.Header("key", "sudis.master.console")
	req.Body = []byte(fmt.Sprintf("%x", md5.Sum([]byte(timestamp+conf.Config.Master.SecurityToken))))
	return req
}

var shutdownCmd = &cobra.Command{
	Use: "shutdown", Short: "关闭master", Long: "关闭master人服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := rpc.NewClient(conf.Config.Master.Band, rpc.OK, nil)

		if err := client.Start(); err != nil {
			return err
		} else {
			defer client.Shutdown()
			//auth
			{
				if resp := client.Send(authRequest(), time.Second*3); resp.Error != nil {
					return resp.Error
				}
			}
			//shutdown
			{
				req := &rpc.Request{URL: "shutdown"}
				if resp := client.Send(req, time.Second*7); resp.Error != nil {
					return resp.Error
				} else {
					fmt.Println("shutdown result : ", string(resp.Body))
				}
			}
			return nil
		}
	},
}

var initCmd = &cobra.Command{
	Use: "init", Short: "初始化数据库", Long: "初始化数据库，创建所需要的表。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dao.InitDatabase(conf.Config.Master.Database)
	},
}

func init() {
	Cmd.AddCommand(shutdownCmd, initCmd)
}
