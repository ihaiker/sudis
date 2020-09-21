package initd

import (
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/sudis/libs/config"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const serviceContent = `
[Unit]
Description=The sudis endpoint.
After=network.target remote-fs.target nss-lookup.target

[Service]
WorkingDirectory=/opt/sudis
ExecStart=/opt/sudis/sudis
ExecStop=/bin/kill -2 $MAINPID
KillSignal=SIGQUIT
TimeoutStopSec=15
KillMode=process
PrivateTmp=true

[Install]
WantedBy=multi-user.target
`

func gt26() bool {
	defer errors.Catch()

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

func writeConfig(confDir string) error {
	configFile := filepath.Join(confDir, "sudis.yaml")
	cfg, _ := yaml.Marshal(config.Config)
	return echo(string(cfg), configFile)
}

func linuxGt26() (err error) {
	defer errors.Catch(func(re error) { err = re })

	errors.Assert(mkdir("/etc/sudis/programs"))
	errors.Assert(writeConfig("/etc/sudis"))

	fileName := "/lib/systemd/system/sudis.service"
	errors.Assert(echo(serviceContent, fileName))

	err, _ = runs("chmod", "+x", fileName)
	errors.Assert(err)

	from, _ := filepath.Abs(os.Args[0])
	to := "/opt/sudis/sudis"
	if from != to {
		_, _ = runs("rm", "-f", to)
		err, _ = runs("cp", "-r", from, to)
		errors.Assert(err)
		err, _ = runs("chmod", "+x", to)
		errors.Assert(err)
	}

	//启动开机启动
	err, out := runs("systemctl", "enable", fileName)
	errors.Assert(err)
	fmt.Println(out)
	return
}

func linuxAutoStart() error {
	if gt26() {
		return linuxGt26()
	} else {
		return errors.New("not support")
	}
}
