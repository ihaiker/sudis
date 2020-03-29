package config

import (
	"github.com/ihaiker/gokit/commons"
	"os"
	"time"
)

type (
	Database struct {
		Type string `mapstructure:"database.type"`
		Url  string `mapstructure:"database.url"`
	}

	sudisConfig struct {
		Debug bool `mapstructure:"debug"`

		//集群唯一标识
		Key string `mapstructure:"key"`

		//绑定地址
		Address string `mapstructure:"address"`

		DisableWebUI bool `mapstructure:"disable-webui"`

		//数据存储位置
		DataPath string `mapstructure:"data-path"`

		Database *Database

		//管理节点的绑定地址
		Manager string `mapstructure:"manager"`

		//节点盐值，如果设置了此值，所有节点的将统一使用此值，如果没有设置，所有节点的盐值都是单独的。
		Salt string `mapstructure:"salt"`

		//连接主节点
		Join []string `mapstructure:"join"`

		//管理程序关闭最大等待时间，防止有些程序不能很快停止而导致的直接kill
		MaxWaitTimeout time.Duration `mapstructure:"maxwait"`

		//时间通知是否同步通知，及只有上一个通知成功后，才可以进行下一个的通知，
		NotifySynchronize bool `mapstructure:"notify-sync"`

		StartTime time.Time `mapstructure:"-"`
	}
)

var Config = defaultConfig()

func autoKey() string {
	name, err := os.Hostname()
	if err != nil || name == "" {
		name = commons.GetHost([]string{"docker0"}, []string{})
	}
	return name
}

func defaultConfig() *sudisConfig {
	return &sudisConfig{
		Key: autoKey(),
		Database: &Database{
			Type: "sqlite3", Url: "sudis.db",
		},
		Join:           []string{},
		MaxWaitTimeout: time.Second * 15,
		StartTime:      time.Now(),
	}
}
