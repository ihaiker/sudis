{
    "msgtype": "markdown",
    "markdown": {
        "title":"Sudis通知",
{{ if eq .Type "process" }}
        "text": "#### 程序状态通知 \n> 节点：{{.Process.Node}}，程序：{{.Process.Name}} \n > {{.FromStatus}} => {{.ToStatus}} \n\n >###### 此通知由sudis自动发送"
{{else}}
        "text": "#### 节点状态通知 \n> {{.Node}} 状态变更为：{{.Status}} \n\n > ###### 此通知由sudis自动发送"
{{end}}
    }
}