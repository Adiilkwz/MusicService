package email

import (
	"fmt"
	"net/smtp"

	"auth_service/internal/domain"
)

type smtpSender struct {
	host     string
	port     string
	from     string
	password string
}

func NewSMTPSender(host, port, from, password string) domain.EmailSender {
	return &smtpSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
	}
}

func (s *smtpSender) SendWelcomeEmail(toEmail, displayName string) error {
	subject := "Welcome to Music Streaming!"
	body := fmt.Sprintf("Hi, %s!\n\nThank you for joining us. Enjoy music!", displayName)

	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) SendPasswordResetEmail(toEmail, resetCode string) error {
	subject := "Passworkd Reset"
	body := fmt.Sprintf("Your code to reser password: %s\n\nIf you did not request for password reset, ignore this message.", resetCode)

	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) sendMail(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"+
		"%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	address := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(address, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
