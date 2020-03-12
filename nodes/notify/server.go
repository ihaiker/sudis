package notify

import (
	"github.com/ihaiker/gokit/concurrent/executors"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/join"
	"github.com/ihaiker/sudis/nodes/notify/mail"
	"github.com/ihaiker/sudis/nodes/notify/webhook"
)

var logger = logs.GetLogger("notity")

type (
	NotifyEventType string

	NotifyEvent struct {
		Type                   NotifyEventType `json:"type"`
		*cluster.NodeEvent     `json:",omitempty"`
		*daemon.FSMStatusEvent `json:",omitempty"`
	}

	notifyServer struct {
		clusterManger *cluster.DaemonManager
		joinManager   *join.ToJoinManager
		executor      executors.ExecutorService
	}
)

const (
	Node    NotifyEventType = "node"
	Process NotifyEventType = "process"
)

func New(sync bool, clusterManger *cluster.DaemonManager, joinManager *join.ToJoinManager) *notifyServer {
	s := &notifyServer{
		clusterManger: clusterManger, joinManager: joinManager,
	}
	if sync {
		s.executor = executors.Single(30 /*任务堆个数*/)
	} else {
		s.executor = executors.Fixed(5 /*工作人数*/, 30 /*任务对个数*/)
	}
	return s
}

func (self *notifyServer) Start() error {
	self.clusterManger.SetStatusListener(func(event daemon.FSMStatusEvent) {
		err := self.executor.Submit(func() {
			self.onProcessStatusEvent(event)
		})
		if err != nil {
			logger.Warn("send process status notify error ", err)
		}
	})
	self.clusterManger.SetNodeListener(func(event cluster.NodeEvent) {
		err := self.executor.Submit(func() {
			self.onNodeStatusEvent(event)
		})
		if err != nil {
			logger.Warn("send node status notify error", err)
		}
	})
	return nil
}

func (self *notifyServer) onProcessStatusEvent(event daemon.FSMStatusEvent) {
	//通知主控节点
	self.joinManager.OnProgramStatusEvent(event)

	//本地处理
	logger.Infof("node: %s, name: %s, from: %s, to: %s",
		event.Process.Node, event.Process.Name, event.FromStatus, event.ToStatus)

	self.notify(NotifyEvent{
		Type:           Process,
		FSMStatusEvent: &event,
	})
}

func (self *notifyServer) onNodeStatusEvent(event cluster.NodeEvent) {
	logger.Infof("node: %s, status: %s", event.Node, event.Status)

	self.notify(NotifyEvent{
		Type:      Node,
		NodeEvent: &event,
	})
}

func (self *notifyServer) notify(data interface{}) {
	webhook.SendWebhook(data)
	mail.SendEmail(data)
}

func (self *notifyServer) Stop() error {
	self.executor.Shutdown()
	return nil
}
