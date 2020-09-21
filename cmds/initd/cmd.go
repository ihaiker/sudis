package initd

import (
	"errors"
	"github.com/spf13/cobra"
	"os/user"
	"runtime"
)

func isAdminUser() error {
	if u, err := user.Current(); err != nil {
		return err
	} else if u.Name == "root" || u.Name == "Administrator" {
		return nil
	} else {
		return errors.New("mast run as root or Administrator")
	}
}

var Cmd = &cobra.Command{
	Use: "initd", Short: "添加开机启动项", Example: "sudis initd",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := isAdminUser(); err != nil {
			return err
		}

		switch runtime.GOOS {
		case "linux":
			return linuxAutoStart()
		case "windows":
			return windowsAutoStart()
		default:
			return errors.New("not support")
		}
	},
}
