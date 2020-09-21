package http

import (
	"encoding/json"
	"github.com/ihaiker/gokit/errors"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/ihaiker/sudis/nodes/notify/mail"
	"github.com/ihaiker/sudis/nodes/notify/webhook"
	"github.com/kataras/iris/v12"
	"time"
)

type NotifyController struct{}

func (self *NotifyController) get(ctx iris.Context) *dao.Notify {
	name := ctx.Params().GetString("name")
	notify, has, err := dao.NotifyDao.Get(name)
	errors.Assert(err)
	errors.True(has, ErrNotFoundConfig)
	return notify
}

func (self *NotifyController) delete(name string) int {
	errors.Assert(dao.NotifyDao.Remove(name), "删除通知配置异常")
	return iris.StatusNoContent
}

func (self *NotifyController) modity(ctx iris.Context) int {
	notify := new(dao.Notify)
	errors.Assert(ctx.ReadJSON(notify))

	errors.Assert(dao.NotifyDao.Remove(notify.Name))

	notify.CreateTime = time.Now()
	errors.Assert(dao.NotifyDao.Add(notify))
	return iris.StatusNoContent
}

func (self *NotifyController) test(ctx iris.Context) int {

	nt := new(dao.Notify)
	errors.Assert(ctx.ReadJSON(nt))

	config := new(dao.JSON)
	errors.Assert(json.Unmarshal([]byte(nt.Config), config))

	if nt.Name == "email" {
		from := config.String("name")
		to := config.String("to")
		content := config.String("content")
		server := mail.NewServer(
			config.String("address"), config.Int("port", 465),
			from, config.String("passwd"),
		)
		errors.Assert(server.Start())
		defer server.Close()

		errors.Assert(server.Send(mail.NewMessage(from, to, "Sudis通知", content)))

	} else {

		errors.Assert(webhook.WebHook(
			config.String("address"),
			config.String("token"),
			[]byte(config.String("content")),
		))

	}
	return iris.StatusNoContent
}
