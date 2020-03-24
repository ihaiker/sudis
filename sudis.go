package main

import (
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/cmds/console"
	"github.com/ihaiker/sudis/cmds/initd"
	"github.com/ihaiker/sudis/cmds/node"
	"github.com/ihaiker/sudis/libs/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	VERSION        = "v3.0.0"
	BUILD_TIME     = "2012-12-12 12:12:12"
	GITLOG_VERSION = "0000"
)

var rootCmd = &cobra.Command{
	Use: filepath.Base(os.Args[0]), Version: fmt.Sprintf(" %s ", VERSION),
	Long: fmt.Sprintf(`SUDIS V3, 一个分布式进程控制程序。
Build: %s, Go: %s, GitLog: %s`, BUILD_TIME, runtime.Version(), GITLOG_VERSION),
	RunE: func(cmd *cobra.Command, args []string) error {
		return node.Start()
	},
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("SUDIS")

		dataPath := viper.GetString("data-path")
		if conf := viper.GetString("conf"); conf != "" {
			viper.SetConfigFile(conf)
		} else if dataPath != "" && files.IsExistFile(filepath.Join(dataPath, "sudis.conf")) {
			viper.SetConfigFile(filepath.Join(dataPath, "sudis.conf"))
		} else {
			viper.SetConfigName("sudis")
			viper.SetConfigType("yaml")
			for _, configPath := range []string{
				"etc", "conf", "etc/sudis",
				"/etc/sudis", os.ExpandEnv("$HOME/.sudis"),
			} {
				viper.AddConfigPath(configPath)
			}
		}
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("read config error: ", err)
		}
		if err := viper.Unmarshal(config.Config); err != nil {
			fmt.Println("unmarshal error ", err)
		}
		if config.Config.DataPath == "" {
			if useConfig := viper.ConfigFileUsed(); useConfig != "" {
				config.Config.DataPath = filepath.Dir(useConfig)
			} else {
				config.Config.DataPath = "$HOME/.sudis"
			}
		}
		config.Config.DataPath = os.ExpandEnv(config.Config.DataPath)
	})

	rootCmd.PersistentFlags().StringP("conf", "f", "", "配置文件")
	rootCmd.PersistentFlags().BoolP("debug", "d", config.Config.Debug, "Debug模式")
	rootCmd.PersistentFlags().StringP("key", "", config.Config.Key, "节点唯一ID")
	rootCmd.PersistentFlags().StringP("address", "", config.Config.Address, "API绑定地址")
	rootCmd.PersistentFlags().BoolP("disable-webui", "", config.Config.DisableWebUI, "禁用webui")

	rootCmd.PersistentFlags().StringP("data-path", "", config.Config.DataPath, "数据存储位置 (default: $HOME/.sudis)")
	rootCmd.PersistentFlags().StringP("database.type", "", config.Config.Database.Type, "数据存储方式")
	rootCmd.PersistentFlags().StringP("database.url", "", config.Config.Database.Url, "数据存储地址")

	rootCmd.PersistentFlags().StringP("salt", "", config.Config.Salt, "安全加密盐值")
	rootCmd.PersistentFlags().String("manager", config.Config.Manager, "管理托管绑定地址")
	rootCmd.PersistentFlags().StringSliceP("join", "", config.Config.Join, "托管连接地址")
	rootCmd.PersistentFlags().DurationP("maxwait", "", config.Config.MaxWaitTimeout, "程序关闭最大等待时间")
	rootCmd.PersistentFlags().BoolP("notify-sync", "", false, "事件通知是否同步通知。")

	rootCmd.AddCommand(console.ConsoleCommands)
	rootCmd.AddCommand(console.Commands...)
	rootCmd.AddCommand(initd.Cmd)
	_ = viper.BindPFlags(rootCmd.PersistentFlags())

	errors.StackFilter = func(frame runtime.Frame) bool {
		return strings.HasPrefix(frame.Function, "github.com/ihaiker/sudis")
	}
}

func main() {
	logs.Open("sudis")
	defer logs.CloseAll()
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
