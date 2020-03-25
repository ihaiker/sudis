package console

import (
	"encoding/json"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/libs/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var modifyCmd = &cobra.Command{
	Use: "modify", Short: "修改程序", Long: "修改被管理的程序", Args: cobra.RangeArgs(1, 2),
	Example: `sudis [console] modify <programName> [jsonfile]
cat jsonfile | sudis [console] modify <programName>`,
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		request := makeRequest(cmd, "modify")

		var content []byte
		if len(args) == 2 {
			if content, err = files.New(args[1]).ToBytes(); err != nil {
				return
			}
		} else {
			if info, _ := os.Stdin.Stat(); info.Size() > 0 {
				if content, err = ioutil.ReadAll(os.Stdin); err != nil {
					return
				}
			}
		}

		if len(content) == 0 {
			return errors.New("no content")
		}

		program := daemon.NewProgram()
		if err = json.Unmarshal(content, program); err != nil {
			return
		}
		program.Name = args[0]
		request.Body, _ = json.Marshal(program)

		sendRequest(cmd, request)
		return nil
	},
}
