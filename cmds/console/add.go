package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var addCmd = &cobra.Command{
	Use: "add", Short: "添加程序管理", Long: "添加被管理的程序", Args: cobra.MinimumNArgs(1),
	Example: "sudis [console] add [@jsonfile|json]...",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "add"
		programs := []string{}
		for _, arg := range args {
			if strings.HasPrefix(arg, "@") {
				f := files.New(arg[1:])
				if !f.IsFile() {
					return errors.New("the file not found :" + arg[1:])
				}
				if fileContent, err := f.ToString(); err != nil {
					return err
				} else {
					programs = append(programs, fileContent)
				}
			} else {
				programs = append(programs, arg)
			}
		}
		request.Body, _ = json.Marshal(programs)
		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
