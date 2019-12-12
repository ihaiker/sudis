package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var startCmd = &cobra.Command{
	Use: "start", Short: "启动管理的程序", Long: "启动管理的某个程序", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] start <programName>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "start"
		request.Body, _ = json.Marshal(args)
		if resp := client.Send(request, time.Second*8); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
