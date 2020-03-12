package mail

import (
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	passwd := os.Getenv("SMTP_PASSWD")
	to := os.Getenv("SMTP_TO")

	server := NewServer(host, 465, user, passwd)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		t.Log(server.Close())
	}()

	m := NewMessage(user, to, "标题", "<h1>aaa</h1>")
	m.From.Name = "异常通知"

	if err := server.Send(m); err != nil {
		t.Fatal(err)
	}
}
