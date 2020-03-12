package dao

import (
	"time"
)

type Notify struct {
	Name       string    `json:"name" yaml:"name" toml:"name" xorm:"varchar(64) notnull pk 'name'"`
	Config     string    `json:"config" yaml:"config" toml:"config" xorm:"config"`
	CreateTime time.Time `json:"createTime" yaml:"createTime" toml:"createTime" xorm:"createTime"`
}

type notifyDao struct {
}

func (self *notifyDao) Get(name string) (notify *Notify, has bool, err error) {
	notify = new(Notify)
	has, err = engine.Where("name = ?", name).Get(notify)
	return
}

func (self *notifyDao) Add(notify *Notify) error {
	_, err := engine.InsertOne(notify)
	return err
}

func (self *notifyDao) Remove(name string) error {
	_, err := engine.Delete(&Notify{Name: name})
	return err
}

var NotifyDao = new(notifyDao)
