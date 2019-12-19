// +build !bindata

package http

import (
	"github.com/kataras/iris"
)

func httpStatic(app *iris.Application) {
	logger.Info("使用 file resources")

	app.Favicon("webui/dist/favicon.ico")

	app.StaticWeb("/", "webui/dist")
	app.StaticWeb("/js", "webui/dist/js")
	app.StaticWeb("/css", "webui/dist/css")
	app.StaticWeb("/img", "webui/dist/img")
	app.StaticWeb("/fonts", "webui/dist/fonts")

	app.RegisterView(iris.HTML("webui/dist", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
}
