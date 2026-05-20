package usecase

import (
	"context"
	"fmt"

	"streaming_service/internal/domain"
)

type likeUsecase struct {
	likeRepo  domain.LikeRepository
	cacheRepo domain.CacheRepository
}

func NewLikeUsecase(likeRepo domain.LikeRepository, cacheRepo domain.CacheRepository) domain.LikeUsecase {
	return &likeUsecase{likeRepo: likeRepo, cacheRepo: cacheRepo}
}

func (u *likeUsecase) LikeSong(ctx context.Context, userID, songID int64) error {
	if u.cacheRepo != nil {
		_ = u.cacheRepo.Set(ctx, likeCacheKey(userID), []int64{}, 1)
	}
	return u.likeRepo.Like(ctx, userID, songID)
}

func (u *likeUsecase) UnlikeSong(ctx context.Context, userID, songID int64) error {
	if u.cacheRepo != nil {
		_ = u.cacheRepo.Set(ctx, likeCacheKey(userID), []int64{}, 1)
	}
	return u.likeRepo.Unlike(ctx, userID, songID)
}

func (u *likeUsecase) GetLikedSongs(ctx context.Context, userID int64, limit, offset int) ([]int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var cached []int64
	if u.cacheRepo != nil {
		if err := u.cacheRepo.Get(ctx, likeCacheKey(userID), &cached); err == nil {
			return cached, nil
		}
	}

	songs, err := u.likeRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	if u.cacheRepo != nil {
		_ = u.cacheRepo.Set(ctx, likeCacheKey(userID), songs, 60)
	}
	return songs, nil
}

func likeCacheKey(userID int64) string {
	return "likes:user:" + fmt.Sprintf("%d", userID)
}
