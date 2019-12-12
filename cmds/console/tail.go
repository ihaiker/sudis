package console

import (
	"fmt"
	"github.com/ihaiker/gokit/remoting/rpc"
	runtimeKit "github.com/ihaiker/gokit/runtime"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var tailCmd = &cobra.Command{
	Use: "tail", Aliases: []string{"tailf"}, Short: "查看日志", Long: "查看程序控制控制态输出日志", Args: cobra.ExactValidArgs(1),
	Example: "sudis [console] tail[f] [-n <num>] <programName,...>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := cmd.CalledAs()
		num, err := cmd.PersistentFlags().GetInt("num")
		if err != nil {
			return err
		}
		request := new(rpc.Request)
		request.Header("num", strconv.Itoa(num))

		if name == "tail" {
			fmt.Println("line")
		} else {
			kill := runtimeKit.NewListener()
			closeSignal := make(chan struct{})
			go func() {
				for {
					select {
					case <-closeSignal:
						return
					case <-time.After(time.Second):
						fmt.Println(time.Now())
					}
				}
			}()
			return kill.WaitWithTimeout(time.Second*3, func() {
				close(closeSignal)
			})
		}
		return nil
	},
}

func init() {
	tailCmd.PersistentFlags().IntP("num", "n", 10, "日志首次条目")
}
