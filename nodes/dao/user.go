package dao

import "github.com/ihaiker/sudis/libs/errors"

type User struct {
	Name   string `json:"name" yaml:"name" toml:"name" xorm:"varchar(32) notnull pk 'name'"`
	Passwd string `json:"-" xorm:"passwd"`
	Time   string `json:"time" yaml:"time" toml:"time" xorm:"time"`
}

type userDao struct {
}

func (dao *userDao) Get(name string) (user *User, has bool, err error) {
	user = new(User)
	has, err = engine.Where("name = ?", name).Get(user)
	return
}

func (dao *userDao) List() (users []*User, err error) {
	users = []*User{}
	err = engine.Find(&users)
	return
}

func (dao *userDao) Remove(name string) error {
	_, err := engine.Delete(&User{Name: name})
	return err
}

func (dao *userDao) ModifyPasswd(name, passwd string) error {
	if _, has, err := dao.Get(name); err != nil {
		return err
	} else if !has {
		return errors.ErrNotFound
	} else {
		_, err = engine.Update(&User{Passwd: passwd}, &User{Name: name})
		return err
	}
}

func (dao *userDao) Insert(user *User) error {
	_, err := engine.InsertOne(user)
	return err
}

var UserDao = new(userDao)
