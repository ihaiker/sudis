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

type CheckHealth struct {
	CheckAddress string    `json:"url" yaml:"url" yaml:"url"`
	CheckMode    CheckMode `json:"type" yaml:"type" yaml:"type"`
	CheckTtl     int       `json:"ttl" yaml:"ttl" yaml:"ttl"`
	//访问安全token定义，这里面是要定义key=value的，这样兼容性更高一些
	SecretToken string `json:"securityKey,omitempty" yaml:"securityKey,omitempty" yaml:"securityKey,omitempty"`
}

//程序命令，（启动或停止）
type Command struct {
	//程序运行体
	Command string `json:"command" yaml:"command" yaml:"command"`

	//启动参数
	Args []string `json:"args,omitempty" yaml:"args,omitempty" yaml:"args,omitempty"`

	//监控检车接口
	CheckHealth *CheckHealth `json:"health,omitempty" yaml:"health,omitempty" yaml:"health,omitempty"`
}

//监控程序
type Program struct {
	//程序唯一性ID，使用UUID方式
	Id uint64 `json:"id" yaml:"id" yaml:"id"`

	Daemon string `json:"daemon" yaml:"daemon" toml:"daemon"`

	//程序名称
	Name string `json:"name" yaml:"name" yaml:"name"`

	//工作目录
	WorkDir string `json:"workDir" yaml:"workDir" yaml:"workDir"`

	//启动使用用户
	User string `json:"user" yaml:"user" yaml:"user"`

	//环境参数变量
	Envs []string `json:"envs,omitempty" yaml:"envs,omitempty" yaml:"envs,omitempty"`

	//是不是守护程序，如果是需要提供启动和停止命令 前台程序
	Start *Command `json:"start" yaml:"start" toml:"start"`

	//启动停止命令
	Stop *Command `json:"stop,omitempty" yaml:"stop,omitempty" toml:"stop,omitempty"`

	//忽略,deamon类型的程序已经启动，也会直接加入管理
	IgnoreAlreadyStarted bool `json:"ignoreStarted" yaml:"ignoreStarted" yaml:"ignoreStarted"`

	//是否自动启动
	AutoStart bool `json:"autoStart" yaml:"autoStart" toml:"autoStart"`

	//启动周期
	StartDuration int `json:"startDuration" yaml:"startDuration" yaml:"startDuration"`

	//启动重试次数
	StartRetries int `json:"startRetries" yaml:"startRetries" yaml:"startRetries"`

	StopSign syscall.Signal `json:"stopSign" yaml:"stopSign" yaml:"stopSign"`

	//结束运行超时时间
	StopTimeout int `json:"stopTimeout" yaml:"stopTimeout" yaml:"stopTimeout"`

	AddTime    time.Time `json:"addTime" yaml:"addTime" yaml:"addTime"`
	UpdateTime time.Time `json:"updateTime" yaml:"updateTime" yaml:"updateTime"`

	//日志文件位置
	Logger string `json:"logger" yaml:"logger" toml:"logger"`
}

func (this *Program) IsForeground() bool {
	return this.Daemon == "0"
}

func (this *Program) JSON() string {
	bs, _ := json.Marshal(this)
	return string(bs)
}

func NewProgram() *Program {
	currentUser, _ := user.Current()
	return &Program{
		Daemon:        "0",
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
