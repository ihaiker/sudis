package console

import (
	"fmt"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/ihaiker/sudis/libs/config"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

const (
	SUBSCRIBE   = "true"
	UNSUBSCRIBE = "false"
)

var tailCmd = &cobra.Command{
	Use: "tail", Aliases: []string{"tailf"}, Short: "查看日志", Long: "查看程序控制控制态输出日志", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] tail[f] [-n <num>] <programName,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {

		name := args[0]
		subscribeId, _ := uuid.NewV4() //channel id
		subId := strings.ReplaceAll(subscribeId.String(), "-", "")

		//订阅成功就可以，在启动客户端连接的试试已经做好了日志的打印
		request := makeRequest(cmd, "tail", name, SUBSCRIBE, subId)
		request.Header("num", strconv.Itoa(viper.GetInt("num")))
		if response := sendRequest(cmd, request); response.Error != nil {
			fmt.Println(response.Error)
			return nil
		}

		kill := runtimeKit.NewListener()
		kill.AddStop(func() error {
			resp := sendRequest(cmd, makeRequest(cmd, "tail", name, UNSUBSCRIBE, subId))
			return resp.Error
		})
		return kill.WaitTimeout(config.Config.MaxWaitTimeout)
	},
}

func init() {
	tailCmd.PersistentFlags().IntP("num", "n", 10, "日志首次条目")
}
