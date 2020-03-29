package auth

import (
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris"
)

type (
	LoginToken struct {
		Token string `json:"token"`
	}

	Service interface {
		Login(data *dao.JSON) *LoginToken
		Check(ctx iris.Context)
	}
)

func NewService() Service {
	return &basicService{}
}
