package http

import (
	"errors"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/server"
	"github.com/kataras/iris"
	"syscall"
	"time"
)

type ProgramForm struct {
	Node string `json:"node"`
	*daemon.Program
}

func newProgramForm() *ProgramForm {
	return &ProgramForm{Program: daemon.NewProgram()}
}

func (self *ProgramForm) Check() (err error) {
	if self.Node == "" {
		return errors.New("请选择您添加程序的节点")
	}
	if self.Name == "" {
		return errors.New("程序名称不能为空！且必须是字母、数字下滑线组合")
	}
	if !self.IsForeground() {
		if self.Stop == nil || self.Stop.Command == "" {
			return errors.New("后台运行程序必须提供停止命令")
		}
		if self.Start.CheckHealth == nil || self.Start.CheckHealth.CheckAddress == "" {
			return errors.New("后台程序必须提供健康检查方式")
		}
	}
	self.StopSign = syscall.SIGTERM
	self.AddTime = time.Now()
	self.UpdateTime = time.Now()
	return
}

func (self *ProgramForm) From(node string, program *daemon.Program) {
	self.Program = program
	self.Node = node
	if program.IsForeground() {
		self.Daemon = "0"
	} else {
		self.Daemon = "1"
	}
}

type ProgramController struct {
	api server.Api
}

func (self *ProgramController) queryPrograms(ctx iris.Context) []*dao.Program {
	name := ctx.URLParam("name")
	node := ctx.URLParam("node")
	tag := ctx.URLParam("tag")
	status := ctx.URLParam("status")
	page := ctx.URLParamIntDefault("page", 1)
	limit := ctx.URLParamIntDefault("limit", 10)
	programs, err := dao.ProgramDao.List(name, node, tag, status, page, limit)
	AssertErr(err)
	return programs
}

func (self *ProgramController) modifyProgramTag(ctx iris.Context) int {
	json := &JSON{}
	AssertErr(ctx.ReadJSON(json))
	name := json.String("name")
	node := json.String("node")
	tag := json.String("tag")
	add := json.Int("add", 1)
	Assert(tag != "", ErrTagEmpty)

	err := dao.ProgramDao.ModifyTag(name, node, tag, add == 1)
	AssertErr(err)
	return iris.StatusNoContent
}

func (self *ProgramController) addOrModifyProgram(ctx iris.Context) int {
	form := newProgramForm()
	AssertErr(ctx.ReadJSON(form))
	{
		err := form.Check()
		AssertErr(err)
		if _, err := self.api.Get(form.Node, form.Name); err != nil && err.Error() == daemon.ErrNotFound.Error() {
			err = self.api.Add(form.Node, form.Program)
			AssertErr(err)
		} else if err != nil {
			AssertErr(err)
		} else {
			err = self.api.Modify(form.Node, form.Name, form.Program)
			AssertErr(err)
		}
	}
	return iris.StatusNoContent
}

func (self *ProgramController) commandProgram(ctx iris.Context) int {
	json := &JSON{}
	AssertErr(ctx.ReadJSON(json))

	name := json.String("name")
	node := json.String("node")
	command := json.String("command")

	logger.Debug("program command: ", name, " ", node, " ", command)

	switch command {
	case "delete":
		AssertErr(self.api.Remove(node, name))
	case "start":
		AssertErr(self.api.Start(node, name))
	case "stop":
		AssertErr(self.api.Stop(node, name))
	case "restart":
		AssertErr(self.api.Stop(node, name))
		AssertErr(self.api.Start(node, name))
	}
	logger.Debug("program command: ", name, " ", node, " ", command, " success")
	return iris.StatusNoContent
}

func (self ProgramController) programDetail(ctx iris.Context) *ProgramForm {
	name := ctx.URLParam("name")
	node := ctx.URLParam("node")
	process, err := self.api.Get(node, name)
	AssertErr(err)

	return &ProgramForm{
		Node: node, Program: process.Program,
	}
}
