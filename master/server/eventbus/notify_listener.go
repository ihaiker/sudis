package eventbus

import (
	"bytes"
	"encoding/json"
	"github.com/ihaiker/sudis/daemon"
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/notify"
	"os"
	"text/template"
)

type NotifyListener struct {
}

func (self *NotifyListener) OnEvent(event *Event) {
	switch event.Name {
	case "NewNode":
	case "LostNode":
	case "ProgramStatus":
		args := event.Value.(*ProgramStatusEvent)
		if args.NewStatus == "" { //添加
			self.notify(args.Key, args.Name, "添加")
		} else if args.OldStatus == "" { //删除
			self.notify(args.Key, args.Name, "删除")
		} else if args.NewStatus == daemon.Running.String() ||
			args.NewStatus == daemon.Stoped.String() || args.NewStatus == daemon.Fail.String() {
			self.notify(args.Key, args.Name, args.NewStatus)
		}
	}
}

func (self *NotifyListener) notify(node, name, state string) {
	self.webhook(node, name, state)
	self.email(node, name, state)
}

func (self *NotifyListener) emailServer() (server *notify.MailServer, config *dao.JSON, err error) {
	var email *dao.Notify
	var has bool
	if email, has, err = dao.NotifyDao.Get("email"); err != nil || !has {
		err = os.ErrNotExist
		return
	}
	config = new(dao.JSON)
	if err = json.Unmarshal([]byte(email.Config), config); err != nil {
		return
	}
	server = notify.NewServer(
		config.String("address"), config.Int("port", 465),
		config.String("name"), config.String("passwd"),
	)
	err = server.Start()
	return
}

func (self *NotifyListener) template(data *dao.JSON, content string) (string, error) {
	out := bytes.NewBuffer([]byte{})
	if t, err := template.New("master").Parse(content); err != nil {
		return "", err
	} else if err = t.Execute(out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (self *NotifyListener) webhook(node, name, state string) {
	cfg, has, err := dao.NotifyDao.Get("webhook")
	if err != nil || !has {
		logger.Warn("WebHook获取异常：has=", has, " error=", err)
		return
	}
	config := new(dao.JSON)
	if err = json.Unmarshal([]byte(cfg.Config), config); err != nil {
		logger.Warn("WebHook配置异常: ", err)
		return
	}
	address := config.String("address")
	token := config.String("token")
	content := config.String("content")

	if body, err := self.template(&dao.JSON{"Node": node, "Name": name, "State": state}, content); err != nil {
		logger.Warn("WebHook模板错误：", err)
	} else if err = notify.WebHook(address, token, body); err != nil {
		logger.Warn("WebHook通知异常：", err)
	}
}

func (self *NotifyListener) email(node, name, state string) {
	server, config, err := self.emailServer()
	if err != nil {
		logger.Error("获取邮件设置错误：", err)
		return
	}
	defer func() { _ = server.Close() }()

	from := config.String("name")
	to := config.String("to")
	content := config.String("content")

	if message, err := self.template(&dao.JSON{"Node": node, "Name": name, "State": state}, content); err != nil {
		logger.Info("模板错误", err)
	} else {
		msg := notify.NewMessage(from, to, "Sudis通知", message)
		if err = server.Send(msg); err != nil {
			logger.Error("发送邮件：", err)
		}
	}
}
