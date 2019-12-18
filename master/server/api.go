package server

import (
	"errors"
	"github.com/ihaiker/sudis/daemon"
)

var ErrNodeOutline = errors.New("节点不在线")

type Api interface {
	Online(node string) bool
	ListProgramNames(node string) ([]string, error)
	ListProgram(node string) ([]*daemon.Process, error)
	Start(node, name string) error
	Stop(node, name string) error
	Add(node string, program *daemon.Program) error
	Remove(node, name string) error
	Get(node, name string) (*daemon.Process, error)
	Modify(node, name string, program *daemon.Program) error

	TailLogger(node, name, id string, num int, logger daemon.TailLogger) error
	UnTailLogger(node, name, id string) error
}

type ApiWrapper struct {
	Apis []Api
}

func NewApiWrapper() *ApiWrapper {
	return &ApiWrapper{Apis: []Api{}}
}

func (self *ApiWrapper) AddApi(api ...Api) {
	self.Apis = append(self.Apis, api...)
}

func (self *ApiWrapper) Online(node string) bool {
	for _, api := range self.Apis {
		if api.Online(node) {
			return true
		}
	}
	return false
}

func (self *ApiWrapper) ListProgramNames(node string) ([]string, error) {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.ListProgramNames(node)
		}
	}
	return nil, ErrNodeOutline
}
func (self *ApiWrapper) ListProgram(node string) ([]*daemon.Process, error) {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.ListProgram(node)
		}
	}
	return nil, ErrNodeOutline
}

func (self *ApiWrapper) Start(node, name string) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Start(node, name)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) Stop(node, name string) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Stop(node, name)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) Add(node string, program *daemon.Program) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Add(node, program)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) Remove(node, name string) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Remove(node, name)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) Get(node, name string) (*daemon.Process, error) {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Get(node, name)
		}
	}
	return nil, ErrNodeOutline
}

func (self *ApiWrapper) Modify(node, name string, program *daemon.Program) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.Modify(node, name, program)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) TailLogger(node, name, id string, num int, consumer daemon.TailLogger) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.TailLogger(node, name, id, num, consumer)
		}
	}
	return ErrNodeOutline
}

func (self *ApiWrapper) UnTailLogger(node, name, id string) error {
	for _, api := range self.Apis {
		if api.Online(node) {
			return api.UnTailLogger(node, name, id)
		}
	}
	return ErrNodeOutline
}
