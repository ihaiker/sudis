package mail

import (
	"bytes"
	"encoding/json"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/nodes/dao"
	"html/template"
	"os"
)

var logger = logs.GetLogger("email")

func loadEmailServer() (server *MailServer, config *dao.JSON, err error) {
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
	server = NewServer(
		config.String("address"), config.Int("port", 465),
		config.String("name"), config.String("passwd"),
	)
	err = server.Start()
	return
}

func SendEmail(data interface{}) {
	server, config, err := loadEmailServer()
	if err == os.ErrNotExist {
		return
	} else if err != nil {
		logger.Error("获取邮件设置错误：", err)
		return
	}
	defer func() { _ = server.Close() }()

	from := config.String("name")
	to := config.String("to")
	content := config.String("content")

	if message, err := render(data, content); err != nil {
		logger.Warn("邮件模板错误: ", err)
	} else {
		msg := NewMessage(from, to, "Sudis Notify", message)
		if err = server.Send(msg); err != nil {
			logger.Warn("发送邮件：", err)
		}
	}
}

func render(data interface{}, content string) (string, error) {
	out := bytes.NewBuffer([]byte{})
	if t, err := template.New("master").Parse(content); err != nil {
		return "", err
	} else if err = t.Execute(out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}
