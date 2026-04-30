package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"auth_service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type accessUsecase struct {
	repo        domain.UserRepository
	jwtSecret   []byte
	emailSender domain.EmailSender
}

func NewAccessUsecase(repo domain.UserRepository, secret string, emailSender domain.EmailSender) domain.AccessUsecase {
	return &accessUsecase{
		repo:        repo,
		jwtSecret:   []byte(secret),
		emailSender: emailSender,
	}
}

func (u *accessUsecase) Register(ctx context.Context, email, password, displayName string) (int64, error) {
	if email == "" || password == "" {
		return 0, errors.New("email and password are required")
	}

	existingUser, err := u.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return 0, errors.New("user with this email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedBytes),
		DisplayName:  displayName,
	}

	userID, err := u.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	go func() {
		err := u.emailSender.SendWelcomeEmail(email, displayName)
		if err != nil {
			fmt.Printf("Error sending welcome email to %s: %v\n", email, err)
		}
	}()

	return userID, nil
}

func (u *accessUsecase) Login(ctx context.Context, email, password string) (*domain.Session, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := u.generateJWT(user.ID, user.Role, time.Minute*15)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.generateJWT(user.ID, user.Role, time.Hour*24*7)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Role:         user.Role,
	}, nil
}

func (u *accessUsecase) ValidateToken(ctx context.Context, accessToken string) (int64, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return u.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid user_id in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("invalid role in token")
	}

	return int64(userIDFloat), role, nil
}

func (u *accessUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	userID, _, err := u.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found or banned")
	}

	newAccessToken, err := u.generateJWT(user.ID, user.Role, time.Minute*15)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := u.generateJWT(user.ID, user.Role, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       user.ID,
		Role:         user.Role,
	}, nil
}

func (u *accessUsecase) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.PasswordHash = string(hashedBytes)
	return u.repo.Update(ctx, user)
}

func (u *accessUsecase) SendPasswordReset(ctx context.Context, email string) error {
	_, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return errors.New("failed to generate reset code")
	}
	resetCode := hex.EncodeToString(bytes)

	err = u.repo.SetResetCode(ctx, email, resetCode)
	if err != nil {
		return errors.New("failed to save reset code")
	}

	go func() {
		err := u.emailSender.SendPasswordResetEmail(email, resetCode)
		if err != nil {
			fmt.Printf("Error sending reset email to %s: %v\n", email, err)
		}
	}()

	return nil
}

func (u *accessUsecase) ConfirmPasswordReset(ctx context.Context, email, resetCode, newPassword string) error {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("invalid request")
	}

	if user.ResetCode == "" || user.ResetCode != resetCode {
		return errors.New("invalid or expired reset code")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.PasswordHash = string(hashedBytes)
	user.ResetCode = ""

	return u.repo.Update(ctx, user)
}

func (u *accessUsecase) generateJWT(userID int64, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}
