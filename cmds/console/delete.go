package console

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var deleteCmd = &cobra.Command{
	Use: "delete", Aliases: []string{"remove"}, Short: "删除管理的程序", Long: "删除被管理的程序", Args: cobra.MinimumNArgs(1),
	Example: "sudis [console] delete <programName,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "delete"
		request.Body, _ = json.Marshal(args)

		if skip, err := cmd.PersistentFlags().GetBool("skip"); err != nil {
			return err
		} else {
			request.Header("skip", strconv.FormatBool(skip))
		}

		if resp := client.Send(request, time.Minute*time.Duration(len(args))); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}

func init() {
	deleteCmd.PersistentFlags().BoolP("skip", "", false, "不停止程序删除")
}
