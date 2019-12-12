package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var modifyCmd = &cobra.Command{
	Use: "modify", Short: "修改程序", Long: "修改被管理的程序", Args: cobra.ExactValidArgs(2),
	Example: "sudis [console] modify <programName> <@jsonfile|json>",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		request := new(rpc.Request)
		request.URL = "modify"

		name := args[0]
		file := args[1]
		content := ""
		if strings.HasPrefix(file, "@") {
			f := files.New(file[1:])
			if !f.IsFile() {
				return errors.New("the file not found :" + f.GetPath())
			}
			if fileContent, err := f.ToString(); err != nil {
				return err
			} else {
				content = fileContent
			}
		} else {
			content = file
		}
		program := daemon.NewProgram()
		if err := json.Unmarshal([]byte(content), program); err != nil {
			return err
		}
		program.Name = name
		request.Body, _ = json.Marshal(program)

		if resp := client.Send(request, time.Second*5); resp.Error != nil {
			fmt.Println(resp.Error)
		} else {
			fmt.Println(string(resp.Body))
		}
		return nil
	},
}
