package console

import (
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use: "stop", Short: "停止管理的程序", Long: "停止正在运行的某个程序", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] stop <programName>",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "stop", args...)
		sendRequest(cmd, request)
	},
}
