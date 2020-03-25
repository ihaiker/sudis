package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

var deleteCmd = &cobra.Command{
	Use: "delete", Aliases: []string{"remove"}, Short: "删除管理的程序", Long: "删除被管理的程序", Args: cobra.MinimumNArgs(1),
	Example: "sudis [console] delete <programName,...>",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "delete", args...)
		if viper.GetBool("skip") {
			request.Header("skip", strconv.FormatBool(true))
		}
		sendRequest(cmd, request)
	},
}

func init() {
	deleteCmd.PersistentFlags().BoolP("skip", "", false, "不停止程序删除")
}
