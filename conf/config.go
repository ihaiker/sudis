package conf

type Server struct {
	//客户端唯一标识
	Key string `json:"key" yaml:"key" toml:"key"`

	//数据文件位置
	Dir string `json:"dir" yaml:"dir" toml:"dir"`

	Sock string `json:"sock" yaml:"sock" toml:"sock"`

	//连接服务端程序
	Master string `json:"master" yaml:"master" toml:"master"`

	//连接管理端安全码
	SecurityToken string `json:"securityToken" yaml:"securityToken" toml:"securityToken"`
}

type Database struct {
	Type string `json:"type" yaml:"type" toml:"type"`
	Url  string `json:"url" yaml:"url" toml:"url"`
}

type Master struct {
	//master绑定端口
	Bind string `json:"bind" yaml:"bind" toml:"bind"`

	//HTTP管理接口
	Http string `json:"http" yaml:"http" toml:"http"`

	//启用ws管理协议
	EnableWS bool `json:"enableWs" yaml:"enableWs" toml:"enableWs"`

	//管理端连接使用的安全码
	SecurityToken string `json:"securityToken" yaml:"securityToken" toml:"securityToken"`

	//数据库连接
	Database *Database `json:"database" yaml:"database" toml:"database"`

	//用户加密盐值
	Salt string `json:"salt" yaml:"salt" toml:"salt"`
}

type SudisConfig struct {
	Version string  `json:"-" yaml:"-" toml:"-"`
	Server  *Server `json:"server" yaml:"server" toml:"server"`
	Master  *Master `json:"master" yaml:"master" toml:"master"`
}

var Config *SudisConfig

func init() {
	Config = &SudisConfig{
		Server: &Server{
			Dir:  "/etc/sudis/programs",
			Sock: "unix://etc/sudis/sudis.sock",
		},
		Master: &Master{
			Bind:          "127.0.0.1:5983",
			Http:          "127.0.0.1:5984",
			EnableWS:      false,
			SecurityToken: "4E4AD35C6C0BEB20DC343A1E8F7E32D4",
			Database: &Database{
				Type: "sqlite3",
				Url:  "/etc/sudis/sudis.db",
			},
			Salt: "2CCAKYGBPTCET2S6",
		},
	}
}
