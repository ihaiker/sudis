// +build !bindata

package http

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/kataras/iris"
)

func httpStatic(app *iris.Application) {
	logs.Info("使用 file resources")
	app.Favicon("webui/dist/favicon.ico")
	app.StaticWeb("/static", "webui/dist/static")
	app.RegisterView(iris.HTML("webui/dist", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
}
