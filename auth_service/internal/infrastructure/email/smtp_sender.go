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
	subject := "Добро пожаловать в Music Streaming!"
	body := fmt.Sprintf("Привет, %s!\n\nСпасибо за регистрацию в нашем сервисе. Наслаждайся музыкой!", displayName)

	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) SendPasswordResetEmail(toEmail, resetCode string) error {
	subject := "Восстановление пароля"
	body := fmt.Sprintf("Ваш код для восстановления пароля: %s\n\nЕсли вы не запрашивали сброс пароля, проигнорируйте это письмо.", resetCode)

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
