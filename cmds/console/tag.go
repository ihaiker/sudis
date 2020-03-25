package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tagCommand = &cobra.Command{
	Use: "tag", Short: "添加程序标签", Long: "给程序添加标签", Args: cobra.ExactArgs(2),
	Example: `sudis [console|cli] tag name tag1 
sudis [console|cli] tag --delete name tag1 
`,
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := makeRequest(cmd, "tag", args...)
		if viper.GetBool("delete") {
			request.Header("delete", "true")
		} else if del, err := cmd.PersistentFlags().GetBool("delete"); err == nil && del {
			request.Header("delete", "true")
		}
		sendRequest(cmd, request)
		return nil
	},
}

func init() {
	tagCommand.PersistentFlags().BoolP("delete", "", false, "是否是删除标签")
}
