package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/nodes/dao"
	"text/template"
)

var logger = logs.GetLogger("webhook")

func SendWebhook(data interface{}) {
	cfg, has, err := dao.NotifyDao.Get("webhook")
	if err != nil {
		logger.Warn("WebHook获取异常：", err)
		return
	}
	if !has {
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

	var body []byte
	if content == "" {
		body, err = json.Marshal(data)
	} else {
		body, err = render(data, content)
	}
	if err != nil {
		logger.Warn("发送通知错误：", err)
		return
	}
	if err = WebHook(address, token, body); err != nil {
		logger.Warn("WebHook通知异常：", err)
	}
}

func render(data interface{}, content string) ([]byte, error) {
	out := bytes.NewBuffer([]byte{})
	if t, err := template.New("master").Parse(content); err != nil {
		return nil, err
	} else if err = t.Execute(out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
