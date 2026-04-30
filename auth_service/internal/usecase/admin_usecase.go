package usecase

import (
	"context"
	"errors"
	"fmt"

	"auth_service/internal/domain"
)

type adminUsecase struct {
	repo domain.UserRepository
}

func NewAdminUsecase(repo domain.UserRepository) domain.AdminUsecase {
	return &adminUsecase{
		repo: repo,
	}
}

func (u *adminUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (u *adminUsecase) UpdateUserRole(ctx context.Context, adminID, targetUserID int64, newRole string) error {
	if newRole != "user" && newRole != "admin" && newRole != "artist" {
		return errors.New("invalid role specified")
	}

	adminUser, err := u.repo.GetByID(ctx, adminID)
	if err != nil || adminUser.Role != "admin" {
		return errors.New("forbidden: requires admin privileges")
	}

	targetUser, err := u.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return errors.New("target user not found")
	}

	targetUser.Role = newRole

	err = u.repo.Update(ctx, targetUser)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}
