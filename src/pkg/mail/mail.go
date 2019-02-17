package mail

import (
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/util"
	"github.com/moocss/go-webserver/src/server"
)

type SendMail struct {
	To          []string
	Cc          []string
	Subject     string
	Body        string
}

func (m *SendMail) Send() error {
	return m.send()
}

func (m *SendMail) AsyncSend() {
	go func() {
		m.send()
	}()
}

func (m *SendMail) send() error {
	err := server.Mail.Send(&util.SendMailMessage{
		To:      m.To,
		Cc:      m.Cc,
		Subject: m.Subject,
		Body:    m.Body,
	})

	if err != nil {
		log.Error("send mail failed：", err.Error())
	}
	return err
}
