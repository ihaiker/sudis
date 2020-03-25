package console

import (
	"github.com/spf13/cobra"
)

var detailCmd = &cobra.Command{
	Use: "detail", Aliases: []string{"inspect"},
	Short: "查看配置信息，JSON", Long: "查看某个程序的配置信息，JSON格式", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] detail <programName>",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "detail", args...)
		sendRequest(cmd, request)
	},
}
