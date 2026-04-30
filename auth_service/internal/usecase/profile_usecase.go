package usecase

import (
	"context"
	"errors"
	"fmt"

	"auth_service/internal/domain"
)

type profileUsecase struct {
	repo domain.UserRepository
}

func NewProfileUsecase(repo domain.UserRepository) domain.ProfileUsecase {
	return &profileUsecase{
		repo: repo,
	}
}

func (u *profileUsecase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *profileUsecase) UpdateProfile(ctx context.Context, userID int64, displayName, avatarURL string) error {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (u *profileUsecase) DeleteAccount(ctx context.Context, userID int64) error {
	err := u.repo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}
