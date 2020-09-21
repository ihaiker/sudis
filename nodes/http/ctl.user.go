package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
	"time"
)

type UserController struct{}

func (self *UserController) queryUser(ctx iris.Context) []*dao.User {
	users, err := dao.UserDao.List()
	errors.Assert(err)
	return users
}

func (self *UserController) addUser(ctx iris.Context) int {
	json := &dao.JSON{}
	errors.Assert(ctx.ReadJSON(json))
	user := &dao.User{
		Name:   json.String("name"),
		Passwd: json.String("passwd"),
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}
	errors.Assert(dao.UserDao.Insert(user))
	return iris.StatusNoContent
}

func (self *UserController) deleteUser(ctx iris.Context) int {
	name := ctx.Params().GetString("name")
	err := dao.UserDao.Remove(name)
	errors.Assert(err)
	return iris.StatusNoContent
}

func (self *UserController) modifyPasswd(ctx iris.Context) int {
	json := &dao.JSON{}
	errors.Assert(ctx.ReadJSON(json))
	name := json.String("name")
	passwd := json.String("passwd")
	errors.Assert(dao.UserDao.ModifyPasswd(name, passwd))
	return iris.StatusNoContent
}
