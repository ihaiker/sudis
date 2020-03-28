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

	rootCmd.PersistentFlags().BoolP("debug", "d", config.Config.Debug, "Debug模式")

	rootCmd.AddCommand(node.NodeCommand)
	rootCmd.AddCommand(console.Commands...)
	rootCmd.AddCommand(console.ConsoleCommands)
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
	node.SetDefaultCommand(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
