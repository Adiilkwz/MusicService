package domain

type EmailSender interface {
	SendPasswordResetEmail(email, resetCode string) error
	SendWelcomeEmail(email, displayName string) error
}

type EventPublisher interface {
	PublishUserRegistered(userID int64, email, displayName string) error
}
