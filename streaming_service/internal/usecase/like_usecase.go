package usecase

import (
	"context"

	"streaming_service/internal/domain"
)

type likeUsecase struct {
	likeRepo domain.LikeRepository
}

func NewLikeUsecase(likeRepo domain.LikeRepository) domain.LikeUsecase {
	return &likeUsecase{likeRepo: likeRepo}
}

func (u *likeUsecase) LikeSong(ctx context.Context, userID, songID int64) error {
	return u.likeRepo.Like(ctx, userID, songID)
}

func (u *likeUsecase) UnlikeSong(ctx context.Context, userID, songID int64) error {
	return u.likeRepo.Unlike(ctx, userID, songID)
}

func (u *likeUsecase) GetLikedSongs(ctx context.Context, userID int64, limit, offset int) ([]int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return u.likeRepo.GetByUserID(ctx, userID, limit, offset)
}
