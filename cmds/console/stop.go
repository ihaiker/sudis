package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var stopCmd = &cobra.Command{
	Use: "stop", Short: "停止管理的程序", Long: "停止正在运行的某个程序", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] stop <programName>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "stop"
		request.Body, _ = json.Marshal(args)
		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
