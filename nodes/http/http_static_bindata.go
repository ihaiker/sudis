// +build bindata

package http

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/kataras/iris/v12"
)

func httpStatic(app *iris.Application) {
	logs.Info("使用bindata")

	app.Get("/favicon.ico", func(ctx iris.Context) {
		bs, _ := Asset("webui/dist/favicon.ico")
		_, _ = ctx.Write(bs)
	})
	app.RegisterView(iris.HTML("webui/dist", ".html").Binary(Asset, AssetNames))
	app.StaticEmbedded("/static", "webui/dist/static", Asset, AssetNames)
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
}
