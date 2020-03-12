package initd

import (
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/files"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
}

func windowsAutoStart() error {
	defer errors.Catch()
	//创建文件夹
	dir, _ := filepath.Abs("./conf")
	errors.Assert(writeConfig(dir))

	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	errors.Assert(err)

	exePath, err := filepath.Abs(os.Args[0])
	errors.Assert(err)

	data := &WindowsServerConfig{
		Id: "sudis", Name: "sudis",
		Description: "The sudis endpoint",
		WorkDir:     workDir,
		Start:       exePath, StartArgs: "",
	}

	out, err := files.New(workDir + "/sudis-server.xml").GetWriter(false)
	errors.Assert(err)

	t, err := template.New("master").Parse(windowCfgServer)
	errors.Assert(err)
	errors.Assert(t.Execute(out, data))

	fmt.Println("下载启动服务插件")
	serverExe := files.New(filepath.Join(workDir, "sudis-server.exe"))
	resp, err := http.Get("https://github.com/kohsuke/winsw/releases/download/winsw-v2.3.0/WinSW.NET4.exe")
	errors.Assert(err)

	defer func() { _ = resp.Body.Close() }()
	fw, err := serverExe.GetWriter(false)
	errors.Assert(err)

	defer func() { _ = fw.Close() }()
	_, err = io.Copy(fw, resp.Body)
	errors.Assert(err)

	err, outlines := runs(serverExe.GetPath(), "install")
	errors.Assert(err)
	fmt.Println(outlines)
	return nil
}
