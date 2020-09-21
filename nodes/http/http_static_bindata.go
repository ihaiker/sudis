// +build bindata

package http

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/kataras/iris/v12"
)

func httpStatic(app *iris.Application) {
	logs.Info("使用bindata")
	app.Get("/favicon.ico", func(ctx iris.Context) {
		_, _ = ctx.Write([]byte(_faviconIco))
	})
	app.HandleDir("/", AssetFile())
	app.Get("/", func(ctx iris.Context) {
		//ctx.View("index.html")
		bs, _ := indexHtmlBytes()
		ctx.Write(bs)
	})
}
