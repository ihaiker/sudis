package console

import (
	"errors"
	"github.com/ihaiker/gokit/files"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var addCmd = &cobra.Command{
	Use: "add", Short: "添加程序管理", Long: "添加被管理的程序",
	Example: "sudis [console] add [jsonfile]...",
	PreRunE: preRune, PostRun: runPost,
	RunE: func(cmd *cobra.Command, args []string) error {
		programs := make([]string, 0)
		for _, arg := range args {
			if fileContent, err := files.New(arg).ToString(); err != nil {
				return err
			} else {
				programs = append(programs, fileContent)
			}
		}
		if content, err := ioutil.ReadAll(os.Stdin); err == nil {
			programs = append(programs, string(content))
		}

		if len(programs) == 0 {
			return errors.New("json file not found")
		}

		request := makeRequest("add", programs...)
		sendRequest(request)
		return nil
	},
}
