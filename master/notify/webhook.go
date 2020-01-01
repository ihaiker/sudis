package notify

import (
	"bytes"
	"errors"
	"net/http"
)

func WebHook(address, accessToken string, body string) error {
	resp, err := http.DefaultClient.Post(address+"?access_token="+accessToken, "application/json;charset=UTF-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("请求异常：" + resp.Status)
	}
	return nil
}
