package dao

import (
	"github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/libs/ipapi"
)

type (
	NodeStatus string

	Node struct {
		Tag        string     `json:"tag" xorm:"tag"`
		Key        string     `json:"key" xorm:"varchar(32) notnull pk 'key'"`
		Token      string     `json:"token" xorm:"token"`
		Ip         string     `json:"ip" xorm:"ip"`
		Address    string     `json:"address" xorm:"address"`
		ProgramNum int        `json:"programNum" xorm:"programNum"`
		Status     NodeStatus `json:"status" xorm:"-"`
		Time       string     `json:"time" xorm:"time"`
	}

	nodeDao struct{}
)

const (
	NodeStatusOnline  NodeStatus = "online"
	NodeStatusOutline NodeStatus = "outline"
)

func (self *nodeDao) List() (nodes []*Node, err error) {
	nodes = make([]*Node, 0)
	err = engine.Find(&nodes)
	return
}

//添加新节点token
func (self *nodeDao) ModifyToken(key, token string) error {
	if _, has, err := self.Get(key); err != nil {
		return err
	} else if has {
		_, err := engine.Update(&Node{Token: token, Time: Timestamp()}, &Node{Key: key})
		return err
	} else {
		_, err := engine.InsertOne(&Node{Key: key, Token: token, Time: Timestamp(), Status: NodeStatusOutline})
		return err
	}
}

func (self *nodeDao) Join(ip, key string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		address := ipapi.Get(ip)
		node = &Node{Key: key, Ip: ip, Address: address.String(), Status: NodeStatusOnline, Time: Timestamp()}
		if _, err = engine.InsertOne(node); err != nil {
			return err
		}
	} else {
		if node.Address == "" || node.Ip != ip {
			node.Address = ipapi.Get(ip).String()
		}
		if _, err := engine.Update(&Node{Ip: ip, Address: node.Address, Time: Timestamp(), Status: NodeStatusOnline}, &Node{Key: key}); err != nil {
			return err
		}
	}
	return nil
}

func (self *nodeDao) Remove(key string) error {
	_, err := engine.Delete(&Node{Key: key})
	return err
}

func (self *nodeDao) Lost(ip, key string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if has {
		node.Time = Timestamp()
		node.Status = NodeStatusOutline
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

func (self *nodeDao) ModifyTag(key, tag string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		return errors.ErrNotFound
	} else {
		node.Tag = tag
		if _, err = engine.Cols("tag").Update(&Node{Tag: node.Tag}, &Node{Key: key}); err != nil {
			return err
		}
		return nil
	}
}

var NodeDao = new(nodeDao)
