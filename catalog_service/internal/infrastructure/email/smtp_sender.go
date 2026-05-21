package infrastructure

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type smtpSender struct {
	host     string
	port     string
	from     string
	password string
}

func NewSMTPSender(host, port, from, password string) EmailSender {
	return &smtpSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
	}
}

func (s *smtpSender) SendEmail(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n%s\r\n", to, subject, body))

	address := fmt.Sprintf("%s:%s", s.host, s.port)

	if s.port == "465" {
		tlsConfig := &tls.Config{
			ServerName: s.host,
		}
		conn, err := tls.Dial("tcp", address, tlsConfig)
		if err != nil {
			return fmt.Errorf("tls dial error: %w", err)
		}
		c, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return fmt.Errorf("smtp new client error: %w", err)
		}
		defer c.Quit()
		auth := smtp.PlainAuth("", s.from, s.password, s.host)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
		if err := c.Mail(s.from); err != nil {
			return fmt.Errorf("smtp mail from failed: %w", err)
		}
		if err := c.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt to failed: %w", err)
		}
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("smtp data start failed: %w", err)
		}
		if _, err := w.Write(msg); err != nil {
			w.Close()
			return fmt.Errorf("smtp write failed: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("smtp write close failed: %w", err)
		}
		return nil
	}

	hostOnly := s.host
	c, err := smtp.Dial(address)
	if err != nil {
		auth := smtp.PlainAuth("", s.from, s.password, hostOnly)
		if err2 := smtp.SendMail(address, auth, s.from, []string{to}, msg); err2 != nil {
			return fmt.Errorf("sendmail fallback failed: %w (dial error: %v)", err2, err)
		}
		return nil
	}
	defer c.Quit()
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: hostOnly}
		if err = c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}
	auth := smtp.PlainAuth("", s.from, s.password, hostOnly)
	if ok, _ := c.Extension("AUTH"); ok {
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}
	if err := c.Mail(s.from); err != nil {
		return fmt.Errorf("mail from failed: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("rcpt to failed: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		w.Close()
		return fmt.Errorf("write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}
	return nil
}
