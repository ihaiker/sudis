package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/libs/config"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/ihaiker/sudis/nodes/http/auth"
	"github.com/ihaiker/sudis/nodes/join"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"net"
)

func Routers(app *iris.Application, clusterManger *cluster.DaemonManager, joinManager *join.ToJoinManager) {
	h := hero.New()
	h.Register(func(ctx iris.Context) *dao.JSON {
		data := &dao.JSON{}
		errors.Assert(ctx.ReadJSON(data))
		return data
	})

	authService := auth.NewService()
	app.Post("/login", h.Handler(authService.Login))

	admin := app.Party("/admin", authService.Check)
	{
		admin.Get("/dashboard", h.Handler(dashboard(clusterManger)))

		node := admin.Party("/node")
		{
			ctl := &NodeController{clusterManger: clusterManger}
			node.Post("/token", h.Handler(ctl.addOrModifyNodeToken))
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
			userCtl := &UserController{}
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

			token := ctx.URLParamTrim("token")
			if token == "" {
				token = config.Config.Salt
			}
			errors.True(token == "", ErrToken.Error())

			must := ctx.URLParamTrim("must")
			if must == "true" {
				joinManager.MustJoinIt(address, token)
			} else {
				errors.Assert(joinManager.Join(address, token))
			}
			return iris.StatusNoContent
		}))
	}
}
