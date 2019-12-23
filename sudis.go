package main

import (
	"github.com/ihaiker/gokit/config"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/cmds/console"
	"github.com/ihaiker/sudis/cmds/initd"
	"github.com/ihaiker/sudis/cmds/master"
	"github.com/ihaiker/sudis/cmds/server"
	"github.com/ihaiker/sudis/cmds/single"
	"github.com/ihaiker/sudis/conf"
	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	VERSION        string = "v2.0.0"
	BUILD_TIME     string = "2012-12-12 12:12:12"
	GITLOG_VERSION string = "0000"
)

var rootCmd = &cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Long:    "sudis, 一个分布式进程控制程序。\n\tBuild: " + BUILD_TIME + ", Go: " + runtime.Version() + ", GitLog: " + GITLOG_VERSION,
	Version: VERSION + "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug, err := cmd.Root().PersistentFlags().GetBool("debug"); err != nil {
			return err
		} else if debug {
			logs.SetDebugMode(debug)
		} else if level, err := cmd.Root().PersistentFlags().GetString("level"); err != nil {
			return err
		} else {
			logs.SetLevel(logs.FromString(level))
		}

		//配置文件
		confPath, err := cmd.Root().PersistentFlags().GetString("conf")
		if err != nil {
			return err
		}

		reg := config.NewConfigRegister("sudis", "")
		if confPath != "" {
			reg.AddPath(confPath)
		}
		reg.With(&configor.Config{
			ENVPrefix:  "SUDIS",
			AutoReload: true, AutoReloadInterval: time.Second * 5,
			AutoReloadCallback: func(config interface{}) {
				logs.Info("配置文件自动重新加载")
				conf.Config = config.(*conf.SudisConfig)
			},
		})
		err = reg.Marshal(conf.Config)
		if err != nil {
			return err
		}
		conf.Config.Version = VERSION
		conf.Config.Server.Dir = os.ExpandEnv(conf.Config.Server.Dir)
		conf.Config.Server.Sock = os.ExpandEnv(conf.Config.Server.Sock)
		conf.Config.Master.Database.Url = os.ExpandEnv(conf.Config.Master.Database.Url)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return single.Cmd.RunE(cmd, args)
	},
}

func init() {
	cobra.OnInitialize(func() {})
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug模式")
	rootCmd.PersistentFlags().StringP("level", "l", "info", "日志级别")
	rootCmd.PersistentFlags().StringP("conf", "f", "", "配置文件")
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(master.Cmd)
	rootCmd.AddCommand(console.Cmds...)
	rootCmd.AddCommand(single.Cmd)
	rootCmd.AddCommand(initd.Cmd)
}

func main() {
	defer logs.CloseAll()
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
