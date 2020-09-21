// +build !bindata

package http

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/kataras/iris/v12"
)

func httpStatic(app *iris.Application) {
	logs.Info("使用 file resources")
	app.Favicon("webui/dist/favicon.ico")
	app.HandleDir("/static", iris.Dir("./webui/dist/static"))
	app.RegisterView(iris.HTML("webui/dist", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
}
