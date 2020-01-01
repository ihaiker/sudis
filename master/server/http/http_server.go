package http

import (
	sysContent "context"
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/server"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"runtime"
	"strconv"
	"time"
)

var logger = logs.GetLogger("master")

func getRequestLogs(ctx context.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

func recoverFn(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}
			switch err.(type) {
			case *Error:
				{
					ctx.StatusCode(err.(*Error).Status)
					_, _ = ctx.JSON(err)
				}
			default:
				{
					var stacktrace string
					for i := 2; ; i++ {
						_, f, l, got := runtime.Caller(i)
						if !got {
							break
						}
						stacktrace += fmt.Sprintf("%s:%d\n", f, l)
					}
					logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
					logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
					logMessage += fmt.Sprintf("Trace: %s\n", err)
					logMessage += stacktrace
					ctx.Application().Logger().Error(logMessage)

					ctx.StatusCode(500)

					ctx.ContentType("application/json;charset=UTF-8")
					_, _ = ctx.JSON(&dao.JSON{"error": "InternalServerError", "message": fmt.Sprintf("%v", err)})
				}
			}
			ctx.StopExecution()
		}
	}()
	ctx.Next()
}

var app = newApplication()

func init() {
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		_, _ = ctx.JSON(dao.JSON{"error": "notfound", "message": "the page not found!" + ctx.Request().RequestURI})
	})
	app.Get("/health", func(ctx iris.Context) {
		_, _ = ctx.JSON(dao.JSON{"status": "UP"})
	})
}

func newApplication() *iris.Application {
	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	app.UseGlobal(crs)

	app.Use(recoverFn)
	return app
}

type masterHttpServer struct {
	address         string
	enableManagerWs bool
	api             *server.ApiWrapper
}

//@param address 服务监听地址
//@param enableManagerWs 服务是否启用,websocket地址
func NewMasterHttpServer(address string, enableManagerWs bool, api *server.ApiWrapper) *masterHttpServer {
	return &masterHttpServer{
		address:         address,
		enableManagerWs: enableManagerWs,
		api:             api,
	}
}

func (self *masterHttpServer) Start() error {
	httpStatic(app)
	Routers(self.api, app)
	ec := commons.Async(func() interface{} {
		return app.Run(
			iris.Addr(self.address),
			iris.WithoutBanner,
			iris.WithoutServerError(iris.ErrServerClosed),
		)
	})

	select {
	case err := <-ec:
		if err != nil {
			return err.(error)
		}
	case <-time.After(time.Second):
		//do nothing
	}

	logger.Info("http server start at: ", self.address)
	return nil
}

func (self *masterHttpServer) Stop() error {
	logger.Debug("http server stop.")
	return app.Shutdown(sysContent.TODO())
}
