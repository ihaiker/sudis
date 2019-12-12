package server

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/server"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "server", Short: "守护进程管理端", Long: "启动守护进程管理端",
	RunE: func(cmd *cobra.Command, args []string) error {
		if key, err := cmd.PersistentFlags().GetString("key"); err != nil {
			return err
		} else if key != "" {
			conf.Config.Server.Key = key
		}
		if conf.Config.Server.Key == "" {
			conf.Config.Server.Key = commons.GetHost([]string{"docker0"}, []string{})
		}
		return server.Start()
	},
}

func init() {
	Cmd.PersistentFlags().StringP("key", "k", "", "客户端唯一标识")
}
