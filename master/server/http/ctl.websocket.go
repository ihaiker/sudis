package http

import (
	"github.com/ihaiker/sudis/master/server"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"strings"
)

type LoggerController struct {
	wsServer *websocket.Server
	api      server.Api
}

func NewLoggerController(api server.Api) *LoggerController {
	wsServer := websocket.New(websocket.Config{})

	wsServer.OnConnection(func(client websocket.Connection) {
		name := client.Context().URLParam("name")
		node := client.Context().URLParam("node")
		uid := strings.ReplaceAll(client.ID(), "-", "")
		logger.Debugf("订阅日志：%s,%s,%s ", name, node, uid)
		go func() {
			err := api.TailLogger(node, name, uid, 30, func(id, line string) {
				_ = client.EmitMessage([]byte(line))
			})
			if err != nil {
				_ = client.EmitMessage([]byte(err.Error()))
			}
		}()
		client.OnDisconnect(func() {
			logger.Debugf("取消订阅日志：%s,%s,%s ", name, node, uid)
			if err := api.UnTailLogger(node, name, uid); err != nil {
				logger.Debug("取消订阅logger错误：", name, ",node", node, ",error:", err)
			}
		})
	})
	return &LoggerController{wsServer: wsServer, api: api}
}

func (self *LoggerController) Handler() iris.Handler {
	return self.wsServer.Handler()
}
