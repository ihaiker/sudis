package initd

import (
	"errors"
	"github.com/spf13/cobra"
	"os/user"
	"runtime"
)

var Cmd = &cobra.Command{
	Use: "initd", Short: "添加开启启动项", Long: "开机启动项添加。\nendpoint，可用值：master|server|single",
	Example: "sudis initd <endpoint>", Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactValidArgs(1)(cmd, args); err != nil {
			return err
		}
		if args[0] != "master" && args[0] != "server" && args[0] != "single" {
			return errors.New("endpoint is error. must be master|server|single")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS == "linux" { //check mast root
			if u, err := user.Current(); err != nil {
				return err
			} else if u.Name != "root" {
				return errors.New("mast run as root")
			}
			return linuxAutoStart(args[0])
		} else if runtime.GOOS == "windows" {
			return windowsAutoStart(args[0])
		}
		return errors.New("not support")
	},
}
