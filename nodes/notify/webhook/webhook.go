package webhook

import (
	"bytes"
	"errors"
	"net/http"
)

func WebHook(address, accessToken string, body []byte) error {
	resp, err := http.DefaultClient.Post(address+"?access_token="+accessToken, "application/json;charset=UTF-8", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("请求异常：" + resp.Status)
	}
	return nil
}
