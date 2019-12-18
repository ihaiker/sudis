package console

import (
	"encoding/json"
	"github.com/ihaiker/gokit/remoting/rpc"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var tailCmd = &cobra.Command{
	Use: "tail", Aliases: []string{"tailf"}, Short: "查看日志", Long: "查看程序控制控制态输出日志", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] tail[f] [-n <num>] <programName,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		num, err := cmd.PersistentFlags().GetInt("num")
		if err != nil {
			return err
		}
		id, _ := uuid.NewV4()
		request := new(rpc.Request)
		request.URL = "tail"
		request.Header("num", strconv.Itoa(num))
		request.Body, _ = json.Marshal([]string{name, "true", id.String()})
		if response := client.Send(request, time.Second*5); response.Error != nil {
			return response.Error
		}

		kill := runtimeKit.NewListener()
		return kill.WaitWithTimeout(time.Second*7, func() {
			request := new(rpc.Request)
			request.URL = "tail"
			request.Body, _ = json.Marshal([]string{name, "false", id.String()})
			client.Send(request, time.Second*5)
		})
	},
}

func init() {
	tailCmd.PersistentFlags().IntP("num", "n", 10, "日志首次条目")
}
