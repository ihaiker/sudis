package console

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use: "status", Short: "查看运行状态", Long: "查看某个程序的运行状态", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] status <programName>",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "status", args...)
		sendRequest(cmd, request)
	},
}
