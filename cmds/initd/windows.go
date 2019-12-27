package initd

import (
	"errors"
	"fmt"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/master/dao"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const windowCfgServer = `
<configuration>
    <id>{{.Id}}</id>
    <name>{{.Name}}</name>
	<description>{{.Description}}</description>

    <workingdirectory>{{.WorkDir}}</workingdirectory>
    <priority>Normal</priority>

    <executable>{{.Start}}</executable>
    <arguments>{{.StartArgs}}</arguments>

    <stopexecutable>{{.Stop}}</stopexecutable>
    <stoparguments>{{.StopArgs}}</stoparguments>

    <startmode>Automatic</startmode>
</configuration>
`

type WindowsServerConfig struct {
	Id          string
	Name        string
	Description string
	WorkDir     string
	Start       string
	StartArgs   string
	Stop        string
	StopArgs    string
}

func windowsAutoStart(endpoint string) error {
	//创建文件夹
	dir := files.New("./conf")
	_ = dir.Mkdir()
	if !dir.Exist() {
		return errors.New("创建配置文件夹错误！")
	}
	fmt.Println("创建配置文件夹：", dir.GetPath())

	if err := writeConfig(endpoint, dir.GetPath()); err != nil {
		return err
	}

	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	workDir = strings.ReplaceAll(workDir, files.ListSeparator, "/")
	exePath := workDir + "/" + filepath.Base(os.Args[0])
	data := &WindowsServerConfig{
		WorkDir: workDir,
		Start:   exePath, StartArgs: "",
		Stop: exePath, StopArgs: "",
	}
	if endpoint == "master" {
		data.Id = "sudis-master"
		data.Name = "sudis-master"
		data.Description = "The sudis master endpoint"
		data.StartArgs = " master"
		data.StopArgs = " master shutdown"
	} else if endpoint == "server" {
		data.Id = "sudis-server"
		data.Name = "sudis-server"
		data.Description = "The sudis server endpoint"
		data.StartArgs = " server"
		data.StopArgs = " console shutdown"
	} else {
		data.Id = "sudis"
		data.Name = "sudis"
		data.Description = "The sudis single endpoint"
		data.StartArgs = ""
		data.StopArgs = " shutdown"
	}

	out, err := files.New(workDir + "/windows-server.xml").GetWriter(false)
	if t, err := template.New("master").Parse(windowCfgServer); err != nil {
		return err
	} else if err = t.Execute(out, data); err != nil {
		return err
	}

	serverExe := files.New("./windows-server.exe")
	fmt.Println("下载启动服务插件")
	if resp, err := http.Get("https://github.com/kohsuke/winsw/releases/download/winsw-v2.3.0/WinSW.NET4.exe"); err != nil {
		return err
	} else {
		defer func() {
			_ = resp.Body.Close()
		}()

		if f, err := serverExe.GetWriter(false); err != nil {
			return err
		} else if _, err = io.Copy(f, resp.Body); err != nil {
			return err
		} else if err = f.Close(); err != nil {
			return err
		}
	}

	if endpoint == "master" || endpoint == "single" {
		if err := dao.InitDatabase(conf.Config.Master.Database); err != nil {
			return err
		}
	}

	err, outlines := runs(serverExe.GetPath(), "install")
	fmt.Println(outlines)
	return err
}
