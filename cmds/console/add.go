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
	Example: `sudis [console] add [jsonfile]...
cat jsonfile | sudis [console] add`,
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

		if info, _ := os.Stdin.Stat(); info.Size() > 0 {
			body, _ := ioutil.ReadAll(os.Stdin)
			programs = append(programs, string(body))
		}

		if len(programs) == 0 {
			return errors.New("json file not found")
		}

		request := makeRequest(cmd, "add", programs...)
		sendRequest(cmd, request)
		return nil
	},
}
