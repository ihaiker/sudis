package console

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start", Short: "启动管理的程序", Long: "启动管理的某个程序", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] start <programName>",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "start", args...)
		sendRequest(cmd, request)
	},
}
