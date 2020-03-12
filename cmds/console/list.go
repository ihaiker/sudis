package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use: "list", Short: "查看程序列表", Long: "查看管理程序的列表（名称）", Args: cobra.NoArgs,
	Example: "sudis [console] list --inspect",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest("list")
		if viper.GetBool("inspect") {
			request.Header("inspect", "true")
		}
		sendRequest(request)
	},
}

func init() {
	listCmd.PersistentFlags().BoolP("inspect", "i", false, "显示详情")
}
