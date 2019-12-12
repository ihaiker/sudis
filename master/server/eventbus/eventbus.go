package eventbus

import (
	"github.com/ihaiker/gokit/logs"
)

var logger = logs.GetLogger("master")

type Eventbus struct {
	C         chan *Event
	Listeners []EventListener
}

func (self *Eventbus) Start() error {
	go func() {
		for {
			select {
			case event := <-self.C:
				if event == nil {
					return
				}
				for _, listener := range self.Listeners {
					listener.OnEvent(event)
				}
			}
		}
	}()
	return nil
}

func (self *Eventbus) Stop() error {
	close(self.C)
	return nil
}

func (self *Eventbus) AddListener(lis EventListener) {
	self.Listeners = append(self.Listeners, lis)
}

func (self *Eventbus) Send(event *Event) {
	self.C <- event
}

var Service = &Eventbus{C: make(chan *Event, 10), Listeners: []EventListener{}}

func Send(event *Event) {
	Service.Send(event)
}
