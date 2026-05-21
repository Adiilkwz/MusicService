package usecase

import (
	"fmt"
	"net/smtp"
)

type smtpSender struct {
	host     string
	port     int
	from     string
	password string
}

func NewSMTPSender(host string, port int, from, password string) *smtpSender {
	return &smtpSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
	}
}

func (s *smtpSender) SendPlaylistSharedNotification(toEmail, playlistName, senderName string) error {
	subject := fmt.Sprintf("Playlist Shared: %s", playlistName)
	body := fmt.Sprintf("Hi,\n\n%s has shared the playlist '%s' with you!\n\nCheck it out now!\n\nBest regards,\nMusic Streaming Team", senderName, playlistName)
	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) SendTrendingSongsNotification(toEmail string) error {
	subject := "Check Out Today's Trending Songs!"
	body := "Hi,\n\nNew trending songs are available on our platform!\n\nVisit now to discover what everyone is listening to.\n\nBest regards,\nMusic Streaming Team"
	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) SendMilestoneNotification(toEmail string, songsCount int) error {
	subject := fmt.Sprintf("Milestone: You've Played %d Songs!", songsCount)
	body := fmt.Sprintf("Hi,\n\nCongratulations! You've played %d songs on our platform!\n\nKeep enjoying great music!\n\nBest regards,\nMusic Streaming Team", songsCount)
	return s.sendMail(toEmail, subject, body)
}

func (s *smtpSender) sendMail(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"+
		"%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	err := smtp.SendMail(address, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
