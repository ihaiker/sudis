package join

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/remoting/rpc"
	"github.com/ihaiker/sudis/daemon"
	sudisError "github.com/ihaiker/sudis/libs/errors"
	"math"
	"strings"
	"time"
)

type ToJoinManager struct {
	key, salt    string
	joined       map[string]*joinClient
	shutdown     bool
	OnRpcMessage rpc.OnMessage
}

func New(key, salt string) *ToJoinManager {
	return &ToJoinManager{
		key: key, salt: salt,
		joined: make(map[string]*joinClient),
	}
}

func (self *ToJoinManager) MustJoinIt(address, token string) {
	maxWaitSeconds := 5 * 60
	go func() {
		for i := 0; !self.shutdown; i++ {
			if err := self.Join(address, token); err == nil {
				return
			} else if strings.Contains(err.Error(), sudisError.ErrToken.Code) {
				logger.Warn("Token错误，忽略加入 ", address)
				return
			}
			seconds := int(math.Pow(2, float64(i)))
			if seconds > maxWaitSeconds {
				seconds = maxWaitSeconds
			}
			next := time.Second * time.Duration(seconds)
			logger.Debugf("%s 重试连接主控节点：%s", next.String(), address)
			time.Sleep(next)
		}
	}()
}

func (self *ToJoinManager) Join(address, token string) (err error) {
	//已经连接成功了，这里的操作是为了客户端连接主控节点异常后，
	//使用命令主动再次连接的判断，因为客户端使用了指数递增方式等待，所以后面的等待是时间将会很长
	if _, has := self.joined[address]; has {
		return
	}
	logger.Infof("连接主控节点：%s", address)
	client := newClient(address, token, self.key, self.OnRpcMessage)
	err = client.Start()
	if err != nil {
		logger.Warn("连接主控异常：", err)
		_ = errors.Safe(client.Stop)
		return
	}
	self.joined[address] = client
	return err
}

func (self *ToJoinManager) Leave(address ...string) error {
	if len(address) == 0 {
		for addr, _ := range self.joined {
			address = append(address, addr)
		}
	}
	for _, addr := range address {
		if cli, has := self.joined[addr]; has {
			logger.Infof("to leave join : %s", addr)
			if err := cli.Stop(); err != nil {
				return err
			}
			delete(self.joined, addr)
		} else {
			return fmt.Errorf("leave %s : not found", addr)
		}
	}
	return nil
}

func (self *ToJoinManager) OnProgramStatusEvent(event daemon.FSMStatusEvent) {
	defer errors.Catch()
	request := &rpc.Request{URL: "program.status"}
	request.Body, _ = json.Marshal(&event)
	for _, client := range self.joined {
		client.Notify(request)
	}
}

func (self *ToJoinManager) Start() error {
	return nil
}

func (self *ToJoinManager) Stop() error {
	logger.Info("multi join stop")
	self.shutdown = true
	for _, client := range self.joined {
		_ = errors.Safe(client.Stop)
	}
	return nil
}
