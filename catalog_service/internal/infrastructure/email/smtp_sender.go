package email

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPSender(host string, port int, user, pass, from string) *SMTPSender {
	return &SMTPSender{Host: host, Port: port, Username: user, Password: pass, From: from}
}

func (s *SMTPSender) Send(ctx context.Context, to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	msg := "From: " + s.From + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0; charset=" + "UTF-8" + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	return smtp.SendMail(addr, auth, s.From, to, []byte(msg))
}
