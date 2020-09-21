package http

import (
	sysContent "context"
	"fmt"
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/ihaiker/sudis/nodes/join"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"strconv"
	"strings"
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

func recoverFn(ctx *context.Context) {
	defer func() {
		if rev := recover(); rev != nil {
			if ctx.IsStopped() {
				return
			}
			err := errors.Convert(rev)
			switch tt := errors.Root(err).(type) {
			case *Error:
				{
					ctx.StatusCode(tt.Status)
					_, _ = ctx.JSON(tt)
				}
			default:
				{
					if tt != errors.ErrAssert {
						ctx.Application().Logger().Error(fmt.Sprintf("%v\n%s", err, errors.Stack()))
					}
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

	cache := iris.Cache304(time.Hour * 24 * 30)
	app.UseGlobal(func(ctx iris.Context) {
		if strings.HasPrefix(ctx.Path(), "/static") {
			cache(ctx)
		} else {
			ctx.Next()
		}
	})

	return app
}

type masterHttpServer struct {
	address       string
	webUI         bool
	joinManager   *join.ToJoinManager
	clusterManger *cluster.DaemonManager
}

//@param address 服务监听地址
//@param disableWebUI 服务是否启用,websocket地址
func NewHttpServer(address string, disableWebUI bool, clusterManger *cluster.DaemonManager, joinManager *join.ToJoinManager) *masterHttpServer {
	return &masterHttpServer{
		address: address, webUI: !disableWebUI,
		clusterManger: clusterManger, joinManager: joinManager,
	}
}

func (self *masterHttpServer) Start() error {
	if self.webUI {
		httpStatic(app)
	}

	Routers(app, self.clusterManger, self.joinManager)

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
