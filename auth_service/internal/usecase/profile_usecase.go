package usecase

import (
	"context"
	"fmt"
	"log"

	"auth_service/internal/domain"
)

type profileUsecase struct {
	repo  domain.UserRepository
	cache domain.UserCache
}

func NewProfileUsecase(repo domain.UserRepository, cache domain.UserCache) domain.ProfileUsecase {
	return &profileUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *profileUsecase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	cachedUser, err := u.cache.GetUser(ctx, userID)
	if err == nil && cachedUser != nil {
		log.Printf("[CACHE HIT] Retrieved user %d from Redis", userID)
		return cachedUser, nil
	}

	log.Printf("[CACHE MISS] Fetching user %d from Postgres", userID)
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	go func() {
		if err := u.cache.SetUser(context.Background(), user); err != nil {
			log.Printf("Warning: Failed to cache user %d: %v", userID, err)
		}
	}()

	return user, nil
}

func (u *profileUsecase) UpdateProfile(ctx context.Context, userID int64, displayName, avatarURL string) error {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
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

	go func() {
		u.cache.DeleteUser(context.Background(), userID)
	}()

	return nil
}

func (u *profileUsecase) DeleteAccount(ctx context.Context, userID int64) error {
	err := u.repo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	go func() {
		u.cache.DeleteUser(context.Background(), userID)
	}()

	return nil
}
