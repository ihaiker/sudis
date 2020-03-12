package mail

import (
	"gopkg.in/gomail.v2"
	"time"
)

type Address struct {
	Address string
	Name    string
}

func NewAddress(address, name string) *Address {
	return &Address{
		Address: address,
		Name:    name,
	}
}

type Body struct {
	ContentType string
	Body        string
	Attachments []gomail.PartSetting
}

type Message struct {
	From    *Address   `mail:"From"`
	To      []*Address `mail:"To"`      //接受者
	Cc      []*Address `mail:"Cc"`      //抄送
	Bcc     []*Address `mail:"Bcc"`     //密送
	Subject string     `mail:"Subject"` //标题
	Body    *Body
}

func NewMessage(from, to, subject, body string) *Message {
	return &Message{
		From:    NewAddress(from, ""),
		To:      []*Address{NewAddress(to, "")},
		Subject: subject,
		Body: &Body{
			ContentType: "text/html",
			Body:        body,
			Attachments: nil,
		},
	}
}

func set(m *gomail.Message, header string, address ...*Address) {
	if address != nil {
		gAddress := make([]string, len(address))
		for i, address := range address {
			gAddress[i] = m.FormatAddress(address.Address, address.Name)
		}
		m.SetHeader(header, gAddress...)
	}
}

func (msg Message) Build() *gomail.Message {
	m := gomail.NewMessage()
	set(m, "From", msg.From)
	set(m, "To", msg.To...)
	set(m, "Cc", msg.Cc...)
	set(m, "Bcc", msg.Bcc...)
	m.SetDateHeader("X-Date", time.Now())
	m.SetHeader("Subject", msg.Subject)
	if msg.Body != nil {
		m.SetBody(msg.Body.ContentType, msg.Body.Body, msg.Body.Attachments...)
	}
	return m
}
