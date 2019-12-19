// +build bindata

package http

import (
	"github.com/kataras/iris"
)

func httpStatic(app *iris.Application) {
	logger.Info("使用bindata")

	app.Get("/favicon.ico", func(ctx iris.Context) {
		bs, _ := Asset("webui/dist/favicon.ico")
		_, _ = ctx.Write(bs)
	})
	app.Get("/logo.png", func(ctx iris.Context) {
		bs, _ := Asset("webui/dist/logo.png")
		_, _ = ctx.Write(bs)
	})
	app.Get("/logo2.png", func(ctx iris.Context) {
		bs, _ := Asset("webui/dist/logo2.png")
		_, _ = ctx.Write(bs)
	})
	app.RegisterView(iris.HTML("webui/dist", ".html").Binary(Asset, AssetNames))
	app.StaticEmbedded("/js", "webui/dist/js", Asset, AssetNames)
	app.StaticEmbedded("/css", "webui/dist/css", Asset, AssetNames)
	app.StaticEmbedded("/img", "webui/dist/img", Asset, AssetNames)
	app.StaticEmbedded("/fonts", "webui/dist/fonts", Asset, AssetNames)

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
}
