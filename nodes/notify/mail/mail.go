package mail

import (
	"github.com/ihaiker/gokit/logs"
	"gopkg.in/gomail.v2"
)

type MailServer struct {
	Host     string
	Port     int
	User     string
	Password string
	stmp     gomail.SendCloser
}

func NewServer(host string, port int, user, passwd string) *MailServer {
	return &MailServer{
		Host: host, Port: port,
		User: user, Password: passwd,
	}
}

func (self *MailServer) Start() (err error) {
	m := gomail.NewDialer(self.Host, self.Port, self.User, self.Password)
	self.stmp, err = m.Dial()
	return
}

func (self *MailServer) Send(msg ...*Message) error {
	if self.stmp == nil {
		logs.GetLogger("mail").Debug("Mail service is not configured")
		return nil
	}

	goMails := make([]*gomail.Message, len(msg))
	for i, message := range msg {
		goMails[i] = message.Build()
	}
	return gomail.Send(self.stmp, goMails...)
}

func (self *MailServer) Close() error {
	if self.stmp != nil {
		return self.stmp.Close()
	}
	return nil
}
