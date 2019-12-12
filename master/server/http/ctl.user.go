package http

import (
	"github.com/ihaiker/sudis/master/dao"
	"github.com/kataras/iris"
	"time"
)

type UserController struct{}

func (self *UserController) login(ctx iris.Context) *JSON {
	json := &JSON{}
	AssertErr(ctx.ReadJSON(json))

	name := json.String("name")
	passwd := json.String("passwd")
	user, has, err := dao.UserDao.Get(name)
	AssertErr(err)
	Assert(has, ErrUser)
	Assert(passwd == user.Passwd, ErrUser)
	return generatorAuth(user.Name)
}

func (self *UserController) queryUser(ctx iris.Context) []*dao.User {
	users, err := dao.UserDao.List()
	AssertErr(err)
	return users
}

func (self *UserController) addUser(ctx iris.Context) int {
	json := &JSON{}
	AssertErr(ctx.ReadJSON(json))
	user := &dao.User{
		Name:   json.String("name"),
		Passwd: json.String("passwd"),
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}
	AssertErr(dao.UserDao.Insert(user))
	return iris.StatusNoContent
}

func (self *UserController) deleteUser(ctx iris.Context) int {
	name := ctx.Params().GetString("name")
	err := dao.UserDao.Remove(name)
	AssertErr(err)
	return iris.StatusNoContent
}

func (self *UserController) modifyPasswd(ctx iris.Context) int {
	json := &JSON{}
	AssertErr(ctx.ReadJSON(json))
	name := json.String("name")
	passwd := json.String("passwd")
	AssertErr(dao.UserDao.ModifyPasswd(name, passwd))
	return iris.StatusNoContent
}
