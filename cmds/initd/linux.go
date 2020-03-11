package initd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/master/dao"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const initConfigFile = `
[Unit]
Description=The sudis %s endpoint.
After=network.target remote-fs.target nss-lookup.target

[Service]
ExecStartPre=rm -f /etc/sudis/sudis.sock
ExecStart=/usr/local/bin/sudis %s
ExecStop=/usr/local/bin/sudis %s shutdown
KillSignal=SIGQUIT
TimeoutStopSec=15
KillMode=process
PrivateTmp=true

[Install]
WantedBy=multi-user.target
`

func gt26() bool {
	defer func() { _ = recover() }()

	v26 := semver.MustParse("2.6.0")
	if err, version := runs("uname", "-r"); err == nil {
		idx := strings.Index(version, "-")
		version = version[0:idx]
		if vCheck, err := semver.Parse(version); err == nil {
			return vCheck.GT(v26)
		}
	}
	return true
}

func mkdir(path string) error {
	fmt.Println("创建目录：", path)
	if files.IsExistDir(path) {
		return nil
	}
	if err := mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.Mkdir(path, 0700)
}

func runs(args ...string) (error, string) {

	fmt.Println("运行：", strings.Join(args, " "))

	out := bytes.NewBuffer([]byte{})
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = out
	cmd.Stdout = out
	err := cmd.Run()
	return err, string(out.Bytes())
}

func echo(content, path string) error {
	fmt.Println("输出内容到：", path)
	fmt.Println("<<----------------------------")
	fmt.Println(content)
	fmt.Println("---------------------------->>")
	service := files.New(path)
	if w, err := service.GetWriter(false); err != nil {
		return err
	} else if _, err = w.WriteString(content); err != nil {
		return err
	}
	return nil
}

func writeConfig(endpoint string, confDir string) error {
	confDir = strings.ReplaceAll(confDir, files.ListSeparator, "")

	conf.Config.Server.Dir = confDir + "/programs"
	conf.Config.Server.Sock = "unix:/" + confDir + "/sudis.sock"
	conf.Config.Master.Database.Url = confDir + "/sudis.db"

	if endpoint == "master" {

	} else if endpoint == "server" {

		conf.Config.Server.Key = commons.GetHost([]string{"docker0"}, []string{})

	} else if endpoint == "single" {

		idx := strings.Index(conf.Config.Master.Bind, ":")
		if idx == 0 {
			conf.Config.Server.Master = "tcp://127.0.0.1" + conf.Config.Master.Bind
		} else {
			conf.Config.Server.Master = "tcp://" + conf.Config.Master.Bind
		}
		conf.Config.Server.SecurityToken = conf.Config.Master.SecurityToken

	}

	cfg, _ := json.MarshalIndent(conf.Config, "", "\t")
	if err := echo(string(cfg), confDir+"/sudis.json"); err != nil {
		return err
	}

	return nil
}

func linuxGt26(endpoint string) error {

	if err := mkdir("/etc/sudis/programs"); err != nil {
		return err
	}

	if err := writeConfig(endpoint, "/etc/sudis"); err != nil {
		return err
	}

	fileName := ""
	content := ""
	if endpoint == "master" {
		fileName = "sudis-master.service"
		content = fmt.Sprintf(initConfigFile, endpoint, endpoint, endpoint)
	} else if endpoint == "server" {
		fileName = "sudis-server.service"
		content = fmt.Sprintf(initConfigFile, endpoint, endpoint, "console")
	} else if endpoint == "single" {
		fileName = "sudis.service"
		content = fmt.Sprintf(initConfigFile, endpoint, "", "")
	}

	if err := echo(content, "/lib/systemd/system/"+fileName); err != nil {
		return err
	}

	if err, out := runs("chmod", "+x", "/lib/systemd/system/"+fileName); err != nil {
		return err
	} else {
		fmt.Println(out)
	}

	if endpoint == "master" || endpoint == "single" {
		if err := dao.InitDatabase(conf.Config.Master.Database); err != nil {
			return err
		}
	}

	self, _ := filepath.Abs(os.Args[0])
	toUser := "/usr/local/bin/sudis"
	if self != toUser {
		_, _ = runs("rm", "-f", toUser)
		if err, _ := runs("cp", "-r", self, toUser); err != nil {
			return err
		}
	}

	//启动开机启动
	if err, out := runs("systemctl", "enable", fileName); err != nil {
		return err
	} else {
		fmt.Println(out)
	}
	return nil
}

func linuxAutoStart(endpoint string) error {
	if gt26() {
		return linuxGt26(endpoint)
	} else {
		return errors.New("not support")
	}
}
