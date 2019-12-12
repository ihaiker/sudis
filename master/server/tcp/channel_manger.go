package tcp

import (
	"github.com/ihaiker/gokit/remoting"
	"reflect"
)

type ServerKeyManager struct {
	channelManager remoting.ChannelManager
	keys           map[string]remoting.Channel
}

func NewServerManager() remoting.ChannelManager {
	return &ServerKeyManager{
		channelManager: remoting.NewIpClientManager(),
		keys:           map[string]remoting.Channel{},
	}
}

func (cm ServerKeyManager) Add(channel remoting.Channel) {
	cm.channelManager.Add(channel)
	if key, has := channel.GetAttr("key"); has {
		cm.keys[key.(string)] = channel
	}
}

func isNil(o interface{}) bool {
	if o == nil {
		return true
	}
	return !reflect.ValueOf(o).IsValid()
}

func (cm ServerKeyManager) Get(index interface{}) (channel remoting.Channel, has bool) {
	if isNil(index) {
		return
	}

	if key, match := index.(string); match {
		if channel, has = cm.keys[key]; has {
			return
		}
	}

	channel, has = cm.channelManager.Get(index)
	return
}
func (cm ServerKeyManager) Remove(channel remoting.Channel) {
	key, has := channel.GetAttr("key")
	cm.channelManager.Remove(channel)
	if has {
		delete(cm.keys, key.(string))
	}
}

func (cm ServerKeyManager) Foreach(fn func(channel remoting.Channel)) {
	cm.channelManager.Foreach(fn)
}
