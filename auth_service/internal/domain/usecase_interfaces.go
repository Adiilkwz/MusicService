package domain

import "context"

type AccessUsecase interface {
	Register(ctx context.Context, email, password, displayName string) (int64, error)
	Login(ctx context.Context, email, password string) (*Session, error)
	ValidateToken(ctx context.Context, accessToken string) (int64, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	SendPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, email, resetCode, newPassword string) error
}

type ProfileUsecase interface {
	GetProfile(ctx context.Context, userID int64) (*User, error)
	UpdateProfile(ctx context.Context, userID int64, displayName, avatarURL string) error
	DeleteAccount(ctx context.Context, userID int64) error
}

type AdminUsecase interface {
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	UpdateUserRole(ctx context.Context, adminID, targetUserID int64, newRole string) error
}
