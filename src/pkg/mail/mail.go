package mail

import (
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src"
)

type Mail struct {
	To          []string
	Cc          []string
	Subject     string
	Body        string
}

func (m *Mail) Send() error {
	return m.send()
}

func (m *Mail) AsyncSend() {
	go func() {
		m.send()
	}()
}

func (m *Mail) send() error {
	err := src.Mail.Send(&SendMailMessage{
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
