package auth

import (
	"encoding/base64"
	"github.com/ihaiker/gokit/errors"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
	"strings"
)

type basicService struct{}

func (self *basicService) Login(data *dao.JSON) *LoginToken {
	username := data.String("name")
	password := data.String("passwd")

	self.mustCheck(username, password)
	token := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	return &LoginToken{Token: token}
}

func (self *basicService) mustCheck(username, password string) {
	user, has, err := dao.UserDao.Get(username)
	errors.Assert(err)
	errors.True(has, ErrUser)
	errors.True(password == user.Passwd, ErrUser)
}

func (self *basicService) Check(ctx iris.Context) {
	authorization := ctx.GetHeader("Authorization")
	if strings.HasPrefix(authorization, "Basic ") {
		if out, err := base64.StdEncoding.DecodeString(authorization[6:]); err == nil {
			nameAndValue := strings.SplitN(string(out), ":", 2)
			err := errors.SafeExec(func() { self.mustCheck(nameAndValue[0], nameAndValue[1]) })
			if err == nil {
				ctx.Next()
				return
			}
		}
	}
	ctx.StatusCode(iris.StatusUnauthorized)
}
