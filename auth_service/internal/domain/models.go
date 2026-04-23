package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	DisplayName  string
	AvatarURL    string
	Role         string
	ResetCode    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	AccessToken  string
	RefreshToken string
	UserID       int64
	Role         string
}
