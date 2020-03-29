package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var joinCmd = &cobra.Command{
	Use: "join", Short: "加入主控节点", Long: "将当前节点托管到其他节点管理", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console|cli] join [--must] <address,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {

		request := makeRequest(cmd, "join", args...)
		if viper.GetBool("must") {
			request.Header("must", "true")
		}
		if token := viper.GetString("token"); token != "" {
			request.Header("token", token)
		}
		sendRequest(cmd, request)
		return nil
	},
}

var leaveCmd = &cobra.Command{
	Use: "leave", Short: "离开主节点", Long: "离开主节点，如果存在多个主节点需要明确指定，否则或全部离开",
	Example: "sudis [console|cli] leave [address,...]",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := makeRequest(cmd, "leave", args...)
		if len(args) == 0 {
			request.Body = []byte("[]")
		}
		sendRequest(cmd, request)
		return nil
	},
}

func init() {
	joinCmd.PersistentFlags().StringP("token", "", "", "连接使用的token")
}
