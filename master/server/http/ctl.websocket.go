package http

import (
	"encoding/json"
	"github.com/ihaiker/sudis/master/dao"
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
		client.OnMessage(func(bytes []byte) {
			params := &dao.JSON{}
			if err := json.Unmarshal(bytes, params); err != nil {
				_ = client.EmitMessage([]byte("认证信息错误！！！"))
				return
			}
			user := params.String("user")
			ticket := params.String("ticket")
			name := params.String("name")
			node := params.String("node")
			uid := strings.ReplaceAll(client.ID(), "-", "")

			if generatorAuth(user).String("token") != ticket {
				_ = client.EmitMessage([]byte("认证信息错误！！！"))
				return
			}

			client.OnDisconnect(func() {
				logger.Debugf("取消订阅日志：%s,%s,%s ", name, node, uid)
				if err := api.UnTailLogger(node, name, uid); err != nil {
					logger.Debug("取消订阅logger错误：", name, ",node", node, ",error:", err)
				}
			})

			logger.Debugf("订阅日志：%s,%s,%s ", name, node, uid)
			go func() {
				err := api.TailLogger(node, name, uid, 30, func(id, line string) {
					_ = client.EmitMessage([]byte(line))
				})
				if err != nil {
					_ = client.EmitMessage([]byte(err.Error()))
				}
			}()
		})
	})
	return &LoggerController{wsServer: wsServer, api: api}
}

func (self *LoggerController) Handler() iris.Handler {
	return self.wsServer.Handler()
}
