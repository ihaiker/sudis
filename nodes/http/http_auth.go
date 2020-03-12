package http

import (
	"crypto/md5"
	"fmt"
	"github.com/ihaiker/sudis/libs/config"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"time"
)

func md5auth(slat, username string) string {
	now := time.Now().Format("20060102")
	code := slat + now + username + now + slat
	return fmt.Sprintf("%x", md5.Sum([]byte(code)))
}

//用户登录生成token
func generatorAuth(user string) *dao.JSON {
	slat := config.Config.Salt
	token := md5auth(slat, user)
	return &dao.JSON{"token": token, "user": user}
}

//验证token
func checkAuth(ctx iris.Context) bool {
	user := ctx.GetHeader("x-user")
	token := ctx.GetHeader("x-ticket")
	slat := config.Config.Salt
	tokenOut := md5auth(slat, user)
	return tokenOut == token
}

func authed(ctx context.Context) {
	if !checkAuth(ctx) {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(dao.JSON{
			"error":   "authFail",
			"message": "认证失败",
		})
	} else {
		ctx.Next()
	}
}
