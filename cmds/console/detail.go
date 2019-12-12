package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var detailCmd = &cobra.Command{
	Use: "detail", Short: "查看配置信息，JSON", Long: "查看某个程序的配置信息，JSON格式", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] detail <programName>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "detail"
		request.Body, _ = json.Marshal(args)
		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
