package dao

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
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

func gbk2utf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func getIpAddress(ip string) string {
	defer func() { _ = recover() }()

	if resp, err := http.Get("http://ipaddr.cz88.net/data.php?ip=" + ip); err != nil {
		return err.Error()
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bs, err := ioutil.ReadAll(resp.Body); err != nil {
			return err.Error()
		} else if utf8bs, err := gbk2utf8(bs); err != nil {
			return err.Error()
		} else {
			address := strings.Split(string(utf8bs[11:len(utf8bs)-1]), ",")[1]
			address = strings.ReplaceAll(address, "'", "")
			return address
		}
	}
}

func (self *nodeDao) Add(ip, key string) error {
	if node, has, err := self.Get(key); err != nil {
		return err
	} else if !has {
		address := getIpAddress(ip)
		node = &Node{Key: key, Ip: ip, Address: address, Status: "online", Time: Timestamp()}
		if _, err = engine.InsertOne(node); err != nil {
			return err
		}
	} else {
		if node.Address == "" || node.Ip != ip {
			node.Address = getIpAddress(ip)
		}
		if _, err := engine.Update(&Node{Ip: ip, Address: node.Address, Time: Timestamp(), Status: "online"}, &Node{Key: key}); err != nil {
			return err
		}
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
		return ErrNotExist
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
		return ErrNotExist
	} else {
		node.ProgramNum = num
		if _, err = engine.Update(node, &Node{Key: key}); err != nil {
			return err
		}
		return nil
	}
}

var NodeDao = new(nodeDao)
