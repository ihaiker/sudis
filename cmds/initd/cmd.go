package initd

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "initd", Short: "添加开启启动项", Long: "开机启动项添加。\nendpoint，可用值：master|server|single",
	Example: "sudis initd <endpoint>", Args: cobra.ExactValidArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
