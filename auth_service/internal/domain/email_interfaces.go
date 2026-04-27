package domain

type EmailSender interface {
	SendWelcomeEmail(toEmail, displayName string) error
	SendPasswordResetEmail(toEmail, resetCode string) error
}
