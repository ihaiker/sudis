package http

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/sudis/conf"
	"github.com/ihaiker/sudis/master/server"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/hero"
	"time"
)

func md5auth(slat, username string) string {
	now := time.Now().Format("20060102")
	code := slat + now + username + now + slat
	return fmt.Sprintf("%x", md5.Sum([]byte(code)))
}

//用户登录生成token
func generatorAuth(user string) *JSON {
	slat := conf.Config.Master.Salt
	token := md5auth(slat, user)
	return &JSON{"token": token, "user": user}
}

//验证token
func checkAuth(ctx iris.Context) bool {
	user := ctx.GetHeader("x-user")
	token := ctx.GetHeader("x-ticket")
	slat := conf.Config.Master.Salt
	tokenOut := md5auth(slat, user)
	return tokenOut == token
}

func authed(ctx context.Context) {
	if !checkAuth(ctx) {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(JSON{
			"error":   "authFail",
			"message": "认证失败",
		})
	} else {
		ctx.Next()
	}
}

func Routers(api server.Api, app *iris.Application) {
	h := hero.New()

	userCtl := &UserController{}
	app.Post("/login", h.Handler(userCtl.login))

	admin := app.Party("/admin", authed)
	{
		admin.Get("/dashboard", h.Handler(dashboard))

		node := admin.Party("/node")
		{
			ctl := &NodeController{api: api}
			node.Get("/list", h.Handler(ctl.queryNodeList)) //列表
			node.Post("/tag", h.Handler(ctl.modifyNodeTag)) //打标签
			node.Put("/reload", h.Handler(ctl.forceReload)) //强制同步
		}
		program := admin.Party("/program")
		{
			ctl := &ProgramController{api: api}
			program.Get("/list", h.Handler(ctl.queryPrograms))
			program.Post("/tag", h.Handler(ctl.modifyProgramTag))
			program.Post("/addOrModify", h.Handler(ctl.addOrModifyProgram))
			program.Put("/command", h.Handler(ctl.commandProgram))
			program.Get("/detail", h.Handler(ctl.programDetail))
		}
		wsServer := NewLoggerController(api)
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
		}
	}
}
