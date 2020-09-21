package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/daemon"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
)

type ProgramController struct {
	clusterManger *cluster.DaemonManager
}

func (self *ProgramController) queryPrograms(ctx iris.Context) *dao.Page {
	name := ctx.URLParam("name")
	node := ctx.URLParam("node")
	tag := ctx.URLParam("tag")
	status := ctx.URLParam("status")
	page := ctx.URLParamIntDefault("page", 1)
	limit := ctx.URLParamIntDefault("limit", 12)
	return self.clusterManger.ListPrograms(name, node, tag, status, page, limit)
}

func (self *ProgramController) modifyProgramTag(json *dao.JSON) int {
	name := json.String("name")
	node := json.String("node")
	tag := json.String("tag")
	add := json.Int("add", 1)
	errors.True(tag != "", ErrTagEmpty)
	errors.Assert(self.clusterManger.ModifyProgramTag(name, node, tag, add == 1), "更新失败")
	return iris.StatusNoContent
}

func (self *ProgramController) addOrModifyProgram(ctx iris.Context) int {
	form := daemon.NewProgram()
	errors.Assert(ctx.ReadJSON(form))

	//check
	errors.True(form.Node != "", "请选择您添加程序的节点")
	errors.True(form.Name != "", "程序名称不能为空！且必须是字母、数字下滑线组合")
	if !form.IsForeground() {
		errors.True(form.Stop != nil && form.Stop.Command != "", "后台运行程序必须提供停止命令")
		errors.True(form.Start.CheckHealth != nil && form.Start.CheckHealth.CheckAddress != "", "后台程序必须提供健康检查方式")
	}

	if _, err := self.clusterManger.GetProgram(form.Node, form.Name); err == nil {
		self.clusterManger.MustModifyProgram(form.Node, form.Name, form)
	} else {
		self.clusterManger.MustAddProgram(form.Node, form)
	}

	return iris.StatusNoContent
}

func (self *ProgramController) commandProgram(json *dao.JSON) int {
	name := json.String("name")
	node := json.String("node")
	command := json.String("command")
	logger.Debug("program command: ", name, " ", node, " ", command)
	self.clusterManger.MustCommand(node, name, command)
	logger.Debug("program command: ", name, " ", node, " ", command, " success")
	return iris.StatusNoContent
}

func (self *ProgramController) programDetail(ctx iris.Context) *daemon.Process {
	name := ctx.URLParam("name")
	node := ctx.URLParam("node")
	return self.clusterManger.MustGetProcess(node, name)
}
