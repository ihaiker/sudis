package join

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/gokit/remoting"
	"github.com/ihaiker/gokit/remoting/rpc"
	"math"
	"time"
)

var logger = logs.GetLogger("join")

type joinClient struct {
	client              rpc.RpcClient
	address, key, token string
	shutdown            bool
}

func newClient(address, token, key string, onRpcMessage rpc.OnMessage) *joinClient {
	joinClient := &joinClient{
		address: address, token: token, key: key,
	}
	joinClient.client = rpc.NewClient(address, onRpcMessage, joinClient.reconnect)
	return joinClient
}

func (self *joinClient) reconnect(channel remoting.Channel) {
	if self.shutdown {
		return
	}

	go func() {
		maxWaitSeconds := 5 * 60
		logger.Debug("尝试连接主控节点")
		for i := 0; !self.shutdown; i++ {
			_ = errors.Safe(self.client.Close)

			if err := self.Start(); err == nil {
				logger.Info("重连与TCP主控节点连接成功：", self.address)
				return
			} else {
				logger.Warn("重连主控节点异常：", err)
			}

			seconds := int(math.Pow(2, float64(i)))
			if seconds > maxWaitSeconds {
				seconds = maxWaitSeconds
			}
			time.Sleep(time.Second * time.Duration(seconds))
		}
	}()
}

func (self *joinClient) Notify(req *rpc.Request) {
	self.client.Oneway(req, time.Second*5)
}

func (self *joinClient) authRequest() *rpc.Request {
	req := new(rpc.Request)
	req.URL = "auth"
	timestamp := time.Now().Format("20060102150405")
	req.Header("timestamp", timestamp)
	req.Header("key", self.key)
	code := fmt.Sprintf("%x", md5.Sum([]byte(timestamp+self.token+self.key)))
	req.Body = []byte(code)
	return req
}

func (self *joinClient) Start() (err error) {
	self.shutdown = false
	if err = self.client.Start(); err != nil {
		return
	}
	if resp := self.client.Send(self.authRequest(), time.Second*10); resp.Error != nil {
		err = resp.Error
	}
	return
}

func (self *joinClient) Stop() error {
	self.shutdown = true
	logger.Info("断开主控节点连接")
	return self.client.Close()
}
