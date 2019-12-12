package console

import (
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var shutdownCmd = &cobra.Command{
	Use: "shutdown", Short: "关闭进程管理服务", Long: "关闭进程管理服务",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "shutdown"
		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			return resp.Error
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
