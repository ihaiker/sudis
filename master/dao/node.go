package dao

import (
	"github.com/ihaiker/gokit/commons"
)

type Node struct {
	Tag        string `json:"tag" yaml:"tag" toml:"tag" xorm:"tag"`
	Key        string `json:"key" yaml:"key" toml:"key" xorm:"varchar(32) notnull pk 'key'"`
	Ip         string `json:"ip" yaml:"ip" toml:"ip" xorm:"ip"`
	Address    string `json:"address" yaml:"address" toml:"address" xorm:"address"`
	ProgramNum int    `json:"programNum" yaml:"programNum" toml:"programNum" xorm:"programNum"`
	Status     string `json:"status" yaml:"status" toml:"status" xorm:"status"`
	Time       string `json:"time" yaml:"time" toml:"time" xorm:"time"`
}

type nodeDao struct {
}

func (self *nodeDao) Ready() {
	_, _ = engine.Update(&Node{Status: "ready"}, &Node{})
}

func (self *nodeDao) List() (nodes []*Node, err error) {
	nodes = make([]*Node, 0)
	err = engine.Find(&nodes)
	return
}

func (self *nodeDao) Add(ip, key string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		node = &Node{Key: key, Ip: ip, Status: "online", Time: Timestamp()}
		if _, err = engine.InsertOne(node); err != nil {
			return err
		}
	} else if _, err := engine.Update(&Node{Ip: ip, Time: Timestamp(), Status: "online"}, &Node{Key: key}); err != nil {
		return err
	}
	return nil
}

func (self *nodeDao) Lost(ip, key string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if has {
		node.Time = Timestamp()
		node.Status = "outline"
		if _, err = engine.Update(node, &Node{Key: key}); err != nil {
			return err
		}
	}
	return nil
}

func (self *nodeDao) Get(key string) (node *Node, has bool, err error) {
	node = new(Node)
	has, err = engine.Where("key = ?", key).Limit(1).Get(node)
	return
}

func (self *nodeDao) Modify(key, tag string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		return commons.ErrNotFound
	} else {
		node.Tag = tag
		if _, err = engine.Cols("tag").Update(&Node{Tag: node.Tag}, &Node{Key: key}); err != nil {
			return err
		}
		return nil
	}
}

func (self *nodeDao) UpdateNodesProcessNumber(key string, num int) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		return commons.ErrNotFound
	} else {
		node.ProgramNum = num
		if _, err = engine.Update(node, &Node{Key: key}); err != nil {
			return err
		}
		return nil
	}
}

var NodeDao = new(nodeDao)
