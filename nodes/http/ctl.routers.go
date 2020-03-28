package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/ihaiker/sudis/nodes/join"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"net"
)

func Routers(app *iris.Application, clusterManger *cluster.DaemonManager, joinManager *join.ToJoinManager) {
	h := hero.New()
	h.Register(func(ctx iris.Context) *dao.JSON {
		data := &dao.JSON{}
		errors.Assert(ctx.ReadJSON(data))
		return data
	})

	userCtl := &UserController{}
	app.Post("/login", h.Handler(userCtl.login))

	admin := app.Party("/admin", authed)
	{
		admin.Get("/dashboard", h.Handler(dashboard(clusterManger)))

		node := admin.Party("/node")
		{
			ctl := &NodeController{clusterManger: clusterManger}
			node.Get("/list", h.Handler(ctl.queryNodeList)) //列表
			node.Post("/tag", h.Handler(ctl.modifyNodeTag)) //打标签
			node.Delete("/{key}", h.Handler(ctl.removeNode))
		}

		program := admin.Party("/program")
		{
			ctl := &ProgramController{clusterManger: clusterManger}
			program.Get("/list", h.Handler(ctl.queryPrograms))
			program.Post("/tag", h.Handler(ctl.modifyProgramTag))
			program.Post("/addOrModify", h.Handler(ctl.addOrModifyProgram))
			program.Put("/command", h.Handler(ctl.commandProgram))
			program.Get("/detail", h.Handler(ctl.programDetail))
		}

		wsServer := NewLoggerController(clusterManger)
		app.Any("/admin/program/logs", wsServer.Handler())

		tag := admin.Party("/tag")
		{
			ctl := TagsController{}
			tag.Get("/list", h.Handler(ctl.queryTag))
			tag.Post("/addOrModify", h.Handler(ctl.addOrModify))
			tag.Delete("/{name}", h.Handler(ctl.removeTag))
		}

		user := admin.Party("/user")
		{
			user.Get("/list", h.Handler(userCtl.queryUser))
			user.Post("/add", h.Handler(userCtl.addUser))
			user.Delete("/{name}", h.Handler(userCtl.deleteUser))
			user.Post("/passwd", h.Handler(userCtl.modifyPasswd))
		}

		notify := admin.Party("/notify")
		{
			ctl := new(NotifyController)
			notify.Get("/{name}", h.Handler(ctl.get))
			notify.Post("/test", h.Handler(ctl.test))
			notify.Post("", h.Handler(ctl.modity))
			notify.Delete("/{name}", h.Handler(ctl.delete))
		}

		admin.Any("/join", h.Handler(func(ctx iris.Context) int {
			address := ctx.URLParamTrim("address")
			errors.True(address != "", "加入地址为空")
			_, _, err := net.SplitHostPort(address)
			errors.Assert(err, "加入地址错误")

			must := ctx.URLParamTrim("must")
			if must == "true" {
				joinManager.MustJoinIt(address)
			} else {
				errors.Assert(joinManager.Join(address))
			}
			return iris.StatusNoContent
		}))
	}
}
