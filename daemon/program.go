package daemon

import (
	"encoding/json"
	"os/user"
	"syscall"
	"time"
)

type CheckMode string

const (
	HTTP  CheckMode = "http"
	HTTPS CheckMode = "https"
	TCP   CheckMode = "tcp"
)

type Tags []string

func (tags *Tags) Add(tag string) {
	*tags = append(*tags, tag)
}

func (tags *Tags) Remove(tag string) {
	for i, t := range *tags {
		if t == tag {
			*tags = append((*tags)[:i], (*tags)[i+1:]...)
			break
		}
	}
}

type (
	CheckHealth struct {
		CheckAddress string    `json:"url"`
		CheckMode    CheckMode `json:"type"`
		CheckTtl     int       `json:"ttl"`
		//访问安全token定义，这里面是要定义key=value的，这样兼容性更高一些
		SecretToken string `json:"securityKey,omitempty"`
	}

	//程序命令，（启动或停止）
	Command struct {
		//程序运行体
		Command string `json:"command,omitempty"`

		//启动参数
		Args []string `json:"args,omitempty"`

		//监控检车接口
		CheckHealth *CheckHealth `json:"health,omitempty"`
	}
	//监控程序
	Program struct {
		//程序唯一性ID，使用UUID方式
		Id uint64 `json:"id"`

		Node string `json:"node,omitempty"`

		Daemon string `json:"daemon"`

		//程序名称
		Name string `json:"name"`

		Description string `json:"description,omitempty"`

		//程序标签
		Tags Tags `json:"tags"`

		//工作目录
		WorkDir string `json:"workDir,omitempty"`

		//启动使用用户
		User string `json:"user,omitempty"`

		//环境参数变量
		Envs []string `json:"envs,omitempty"`

		//是不是守护程序，如果是需要提供启动和停止命令 前台程序
		Start *Command `json:"start"`

		//启动停止命令
		Stop *Command `json:"stop,omitempty"`

		//忽略,deamon类型的程序已经启动，也会直接加入管理
		IgnoreAlreadyStarted bool `json:"ignoreStarted,omitempty"`

		//是否自动启动
		AutoStart bool `json:"autoStart,omitempty"`

		//启动周期
		StartDuration int `json:"startDuration,omitempty"`

		//启动重试次数
		StartRetries int `json:"startRetries,omitempty"`

		StopSign syscall.Signal `json:"stopSign,omitempty"`

		//结束运行超时时间
		StopTimeout int `json:"stopTimeout,omitempty"`

		AddTime    time.Time `json:"addTime,omitempty"`
		UpdateTime time.Time `json:"updateTime,omitempty"`

		//日志文件位置
		Logger string `json:"logger,omitempty"`
	}
)

func (this *Program) IsForeground() bool {
	return this.Daemon == "0"
}

func (this *Program) JSON() string {
	bs, _ := json.Marshal(this)
	return string(bs)
}
func (this *Program) JSONByte() []byte {
	bs, _ := json.Marshal(this)
	return bs
}

func NewProgram() *Program {
	currentUser, _ := user.Current()
	return &Program{
		Daemon:        "0",
		Tags:          Tags{},
		WorkDir:       currentUser.HomeDir,
		User:          currentUser.Username,
		AutoStart:     false,
		StartDuration: 7,
		StartRetries:  3,
		Envs:          []string{},
		StopSign:      syscall.SIGQUIT,
		StopTimeout:   7,
		AddTime:       time.Now(),
		UpdateTime:    time.Now(),
	}
}
