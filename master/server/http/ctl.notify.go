package http

import (
	"encoding/json"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/notify"
	"github.com/kataras/iris"
	"time"
)

type NotifyController struct{}

func (self *NotifyController) get(ctx iris.Context) *dao.Notify {
	name := ctx.Params().GetString("name")
	notify, has, err := dao.NotifyDao.Get(name)
	AssertErr(err)
	Assert(has, ErrNotFoundConfig)
	return notify
}

func (self *NotifyController) modity(ctx iris.Context) int {
	notify := new(dao.Notify)
	AssertErr(ctx.ReadJSON(notify))

	AssertErr(dao.NotifyDao.Remove(notify.Name))

	notify.CreateTime = time.Now()
	AssertErr(dao.NotifyDao.Add(notify))
	return iris.StatusNoContent
}

func (self *NotifyController) test(ctx iris.Context) int {

	nt := new(dao.Notify)
	AssertErr(ctx.ReadJSON(nt))

	config := new(JSON)
	AssertErr(json.Unmarshal([]byte(nt.Config), config))

	if nt.Name == "email" {
		from := config.String("name")
		to := config.String("to")
		content := config.String("content")
		server := notify.NewServer(
			config.String("address"), config.Int("port", 465),
			from, config.String("passwd"),
		)
		AssertErr(server.Start())
		defer server.Close()

		AssertErr(server.Send(notify.NewMessage(from, to, "Sudis通知", content)))

	} else {

		AssertErr(notify.WebHook(
			config.String("address"),
			config.String("token"),
			config.String("content"),
		))

	}
	return iris.StatusNoContent
}
