{{ if eq .Type "process" }}
    <h1>程序通知：</h1>
    <table>
        <tr>
            <th>程序节点</th><td>{{.Process.Node}}</td>
            <th>名称</th><td>{{.Process.Name}}</td>
        </tr>
        <tr>
            <th>From:</th><td>{{.FromStatus}}</td>
            <th>To:</th><td>{{.ToStatus}}</td>
        </tr>
    </table>
{{else}}
    <h1>节点通知：</h1>
    {{.Node}} 状态变更为：{{.Status}}
{{end}}