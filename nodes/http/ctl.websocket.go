package http

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

type LoggerController struct {
	wsServer *neffos.Server
	manger   *cluster.DaemonManager
}

func NewLoggerController(manger *cluster.DaemonManager) *LoggerController {
	connectIds := map[string][]string{}

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			params := &dao.JSON{}
			if err := json.Unmarshal(msg.Body, params); err != nil {
				return err
			}
			//token := params.String("token")
			name := params.String("name")
			node := params.String("node")
			line := params.Int("line", 30)
			uid := strings.ReplaceAll(nsConn.Conn.ID(), "-", "")
			logger.Debugf("订阅日志：%s,%s,%s ", name, node, uid)
			if err := manger.SubscribeLogger(uid, node, name, func(id, line string) {
				defer errors.Catch(func(err error) {})
				//nsConn.Emit("", []byte(line))
				fmt.Println(line)
				nsConn.Conn.Write(neffos.Message{Body: []byte(line)})
			}, line); err != nil {
				//nsConn.Emit("", []byte(err.Error()))
				nsConn.Conn.Write(neffos.Message{Body: []byte(err.Error())})
			}
			connectIds[uid] = []string{name, node}
			return nil
		},
	})
	ws.OnDisconnect = func(c *websocket.Conn) {
		uid := strings.ReplaceAll(c.ID(), "-", "")
		if nameAndNode, has := connectIds[uid]; has {
			name, node := nameAndNode[0], nameAndNode[1]
			logger.Debugf("取消订阅日志：%s,%s,%s ", name, node, uid)
			_ = manger.UnsubscribeLogger(uid, node, name)
		}
	}
	return &LoggerController{wsServer: ws, manger: manger}
}

func (self *LoggerController) Handler() iris.Handler {
	return websocket.Handler(self.wsServer)
}
