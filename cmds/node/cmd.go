package node

import (
	"github.com/ihaiker/sudis/libs/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var NodeCommand = &cobra.Command{
	Use: "node", Short: "节点启动",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Start()
	},
}

func init() {
	NodeCommand.PersistentFlags().StringP("conf", "f", "", "配置文件")

	NodeCommand.PersistentFlags().StringP("key", "", config.Config.Key, "节点唯一ID")
	NodeCommand.PersistentFlags().StringP("address", "", config.Config.Address, "API绑定地址")
	NodeCommand.PersistentFlags().BoolP("disable-webui", "", config.Config.DisableWebUI, "禁用webui")

	NodeCommand.PersistentFlags().StringP("data-path", "", config.Config.DataPath, "数据存储位置 (default: $HOME/.sudis)")
	NodeCommand.PersistentFlags().StringP("database.type", "", config.Config.Database.Type, "数据存储方式")
	NodeCommand.PersistentFlags().StringP("database.url", "", config.Config.Database.Url, "数据存储地址")

	NodeCommand.PersistentFlags().StringP("salt", "", config.Config.Salt, "安全加密盐值，如果设置了此值，所有节点加入管理默认将使用此值，若未设置节点将使用单独的设置")
	NodeCommand.PersistentFlags().String("manager", config.Config.Manager, "管理托管绑定地址")
	NodeCommand.PersistentFlags().StringSliceP("join", "", config.Config.Join, "托管连接地址")
	NodeCommand.PersistentFlags().DurationP("maxwait", "", config.Config.MaxWaitTimeout, "程序关闭最大等待时间")
	NodeCommand.PersistentFlags().BoolP("notify-sync", "", false, "事件通知是否同步通知。")

	_ = viper.BindPFlags(NodeCommand.PersistentFlags())
}

func SetDefaultCommand(root *cobra.Command) {
	setDef := NodeCommand
	//set node is default command
	if runCommand, args, err := root.Find(os.Args[1:]); err == nil {
		if runCommand == root {
			root.InitDefaultHelpFlag()
			_ = root.ParseFlags(args)

			if help, err := root.Flags().GetBool("help"); err == nil && help {
				// show help
			} else {
				idx := 1
				for _, arg := range args {
					if strings.HasPrefix(arg, "-") {
						flagName := strings.TrimLeft(arg, "-")
						hasValue := strings.Index(flagName, "=")
						if hasValue > 0 {
							flagName = flagName[:hasValue]
						}
						if f := root.PersistentFlags().Lookup(flagName); f != nil {
							if f.Value.Type() == "bool" || hasValue > 0 {
								idx += 1
							} else if f.Value.String() != "" {
								idx += 2
							}
							continue
						}
						if len(flagName) == 1 {
							if f := root.PersistentFlags().ShorthandLookup(flagName); f != nil {
								if f.Value.Type() == "bool" || hasValue > 0 {
									idx += 1
								} else if f.Value.String() != "" {
									idx += 2
								}
								continue
							}
						}
					}
					break
				}
				root.SetArgs(append(os.Args[1:idx], append([]string{setDef.Name()}, os.Args[idx:]...)...))
			}
		}
	}
}
