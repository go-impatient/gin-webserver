package util

import (
	"gopkg.in/gomail.v2"
)

type SendMail struct {
	Enable  int
	Smtp    string
	Port    int
	User    string
	Pass    string
	dialer  *gomail.Dialer
}

type SendMailMessage struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Body    string
	Attach  string
	mail    *SendMail
}

func SendMailNew(mail *SendMail) *SendMail {
	mail.dialer = gomail.NewPlainDialer(mail.Smtp, mail.Port, mail.User, mail.Pass)
	return mail
}

func (mail *SendMail) Send(msg *SendMailMessage) error {
	if mail.Enable == 0 {
		return nil
	}
	msg.mail = mail
	m := msg.NewMessage()
	if err := mail.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (m *SendMailMessage) NewMessage() *gomail.Message {
	mailMsg := gomail.NewMessage()
	if m.From == "" {
		mailMsg.SetHeader("From", m.mail.User)
	} else {
		mailMsg.SetHeader("From", m.From)
	}
	mailMsg.SetHeader("To", m.To...)
	if len(m.Cc) > 0 {
		mailMsg.SetHeader("Cc", m.Cc...)
	}
	mailMsg.SetHeader("Subject", m.Subject)
	mailMsg.SetBody("text/html", m.Body)
	if m.Attach != "" {
		mailMsg.Attach(m.Attach)
	}
	return mailMsg
}

