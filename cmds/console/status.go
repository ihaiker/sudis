package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var statusCmd = &cobra.Command{
	Use: "status", Short: "查看运行状态", Long: "查看某个程序的运行状态", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] status <programName>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "status"
		request.Body, _ = json.Marshal(args)

		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
