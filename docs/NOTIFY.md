# Sudis 通知配置

> 阅读本章节介绍，您需要先行了解 Golang语言 `html/template`和`text/template`相关内容，前方高能请注意。

sudis 通知只能由webui来填写设置，如果需要设置通知请先打开webui.

sudis的现在支持两种通知，邮件和webhook。当程序状态或者集群节点改变后，直接配置发送给用户。


## 邮件通知

当需要通知时，程序将获取当前邮件配置，如果未配置将不再发送通知。

邮件内容将采用html/tempate渲染用户指定模板并生成邮件内容。

> 邮件模板示例

```
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
```

> 邮件推送模板结构体

```go
type NotifyEvent struct {
		Type string //process/node

		//Type == node
		//{
		Node   string
		Status string  //online/outline
		//}



		//Type == process
		//*{
		Process *struct {
			Pid    int
			Status string //ready/starting/running/fail/retry/stopping/stoped/

			//程序唯一性ID，使用UUID方式
			Id uint64

			Node string

			Daemon string

			//程序名称
			Name string

			Description string

			//程序标签
			Tags []string

			//工作目录
			WorkDir string

			//启动使用用户
			User string

			//环境参数变量
			Envs []string

			//是不是守护程序，如果是需要提供启动和停止命令 前台程序
			Start *struct {
				//程序运行体
				Command string

				//启动参数
				Args []string

				//监控检车接口
				CheckHealth *struct {
					CheckAddress string
					CheckMode    CheckMode
					CheckTtl     int
					SecretToken  string
				}
			}

			//启动停止命令
			Stop  *struct {
				//程序运行体
				Command string

				//启动参数
				Args []string

				//监控检车接口
				CheckHealth *struct {
					CheckAddress string
					CheckMode    CheckMode
					CheckTtl     int
					SecretToken  string
				}
			}

			//忽略,deamon类型的程序已经启动，也会直接加入管理
			IgnoreAlreadyStarted bool

			//是否自动启动
			AutoStart bool

			//启动周期
			StartDuration int

			//启动重试次数
			StartRetries int

			//结束运行超时时间
			StopTimeout int

			AddTime    time.Time
			UpdateTime time.Time

			//日志文件位置
			Logger string

			retryLeft int

			Cpu float64
			Rss uint64
		}
		FromStatus string //ready/starting/running/fail/retry/stopping/stoped/
		ToStatus   string // ready/starting/running/fail/retry/stopping/stoped/
    //}
)

```



## webhook推送

webhook系统采用POST方式推送，body内容采用`text/template`方式处理用户给定模板内容推送。 如果用户为指定特定的模板系统将会使用下面的json发送

> 程序通知

webhook推送可以省略推送模板，如果用户为指定特定的模板系统将会使用下面的`json`发送

```json
{
	"type": "process",
	"process": {
		"pid": 44593,
		"status": "running",
		"id": 2,
		"node": "sudis",
		"daemon": "0",
		"name": "ping",
		"tags": [],
		"workDir": "/Users/haiker",
		"user": "haiker",
		"start": {
			"command": "ping",
			"args": ["172.16.100.2"]
		},
		"stop": {},
		"startDuration": 3,
		"startRetries": 3,
		"stopSign": 3,
		"stopTimeout": 7,
		"addTime": "2020-03-24T12:47:41.025598+08:00",
		"updateTime": "2020-03-24T12:47:41.025598+08:00",
		"cpu": 0,
		"rss": 4187
	},
	"fromStatus": "starting",
	"toStatus": "running"
}
```

> 节点通知

```json
{
  "type":"node",
  "node":"test02",
  "status":"online"
}
```



当然你也可使用模板转换JSON内容。例如钉钉推送的模板

> 钉钉推送模板

```
{
    "msgtype": "markdown",
    "markdown": {
        "title":"Sudis通知",
{{ if eq .Type "process" }}
        "text": "#### 程序状态通知 \n> {{.Process.Node}} {{.Process.Name}}  {{.FromStatus}} => {{.ToStatus}} \n\n >###### 此通知由sudis自动发送"
{{else}}
        "text": "#### 节点状态通知 \n> {{.Node}} 状态变更为：{{.Status}} \n\n > ###### 此通知由sudis自动发送"
{{end}}
    }
}
```

模板使用的数据结构体和可以参见`邮件推送模板结构体`。