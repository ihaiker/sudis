package http

import (
	"encoding/json"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"strings"
)

type LoggerController struct {
	wsServer *websocket.Server
	manger   *cluster.DaemonManager
}

func NewLoggerController(manger *cluster.DaemonManager) *LoggerController {

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
			line := params.Int("line", 30)
			uid := strings.ReplaceAll(client.ID(), "-", "")

			if generatorAuth(user).String("token") != ticket {
				_ = client.EmitMessage([]byte("认证信息错误！！！"))
				return
			}

			client.OnDisconnect(func() {
				logger.Debugf("取消订阅日志：%s,%s,%s ", name, node, uid)
				_ = manger.UnsubscribeLogger(node, name, uid)
			})
			logger.Debugf("订阅日志：%s,%s,%s ", name, node, uid)
			if err := manger.SubscribeLogger(uid, node, name, func(id, line string) {
				defer errors.Catch(func(err error) {})
				_ = client.EmitMessage([]byte(line))
			}, line); err != nil {
				_ = client.EmitMessage([]byte(err.Error()))
			}
		})
	})
	return &LoggerController{wsServer: wsServer, manger: manger}
}

func (self *LoggerController) Handler() iris.Handler {
	return self.wsServer.Handler()
}
