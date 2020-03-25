package console

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use: "list", Aliases: []string{"ls"}, Args: cobra.NoArgs,
	Short: "查看程序列表", Long: "查看管理程序的列表（名称）",
	Example: "sudis [console] list --inspect",
	PreRunE: preRune, PostRun: runPost,
	Run: func(cmd *cobra.Command, args []string) {
		request := makeRequest(cmd, "list")
		if viper.GetBool("inspect") {
			request.Header("inspect", "true")
		}
		if viper.GetBool("all") {
			request.Header("all", "true")
		}
		if viper.GetBool("quiet") {
			request.Header("quiet", "true")
		}

		resp := sendRequest(cmd, request, true)
		if resp.Error != nil {
			fmt.Println(resp.Error)
		} else if viper.GetBool("inspect") {
			fmt.Println(string(resp.Body))
		} else {
			items := make([]string, 0)
			if err := json.Unmarshal(resp.Body, &items); err != nil {
				fmt.Println(err)
			} else {
				for _, item := range items {
					fmt.Println(item)
				}
			}
		}
	},
}

func init() {
	listCmd.PersistentFlags().BoolP("inspect", "i", false, "显示详情")
	listCmd.PersistentFlags().BoolP("all", "a", false, "显示全部程序")
	listCmd.PersistentFlags().BoolP("quiet", "q", false, "仅仅显示名称")
}
