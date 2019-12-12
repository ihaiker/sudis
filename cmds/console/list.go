package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"time"
)

var listCmd = &cobra.Command{
	Use: "list", Short: "查看程序列表", Long: "查看管理程序的列表（名称）",
	Example: "sudis [console] list [inspect]",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "list"
		request.Body, _ = json.Marshal(args)
		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			if len(args) == 0 {
				names := []string{}
				if err := json.Unmarshal(resp.Body, &names); err != nil {
					fmt.Println(err)
				} else {
					for _, name := range names {
						fmt.Println(name)
					}
				}
			} else {
				fmt.Println(string(resp.Body))
			}
		}
		return nil
	},
}
