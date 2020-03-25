package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//TODO leave 方法开发，退出托管

var joinCmd = &cobra.Command{
	Use: "join", Short: "加入主控节点", Long: "将当前节点托管到其他节点管理", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console|cli] join [--must] <address,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {

		request := makeRequest(cmd, "join", args...)
		if viper.GetBool("must") {
			request.Header("must", "true")
		}

		sendRequest(cmd, request)
		return nil
	},
}

func init() {
	joinCmd.PersistentFlags().BoolP("must", "", false, "是否要后台重试连接")
}
