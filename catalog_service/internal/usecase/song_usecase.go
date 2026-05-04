package usecase

import (
	"catalog_service/internal/domain"
	"context"
)

type songUsecase struct {
	songRepo domain.SongRepository
}

func NewSongUsecase(sr domain.SongRepository) domain.SongUsecase {
	return &songUsecase{
		songRepo: sr,
	}
}

func (u *songUsecase) CreateSong(ctx context.Context, albumID int64, title string, duration int32, genre string) (int64, error) {
	song := &domain.Song{
		AlbumID:         albumID,
		Title:           title,
		DurationSeconds: duration,
		Genre:           genre,
	}
	return u.songRepo.Create(ctx, song)
}

func (u *songUsecase) GetSong(ctx context.Context, id int64) (*domain.Song, error) {
	return u.songRepo.GetByID(ctx, id)
}

func (u *songUsecase) UpdateSong(ctx context.Context, id int64, title, genre string) error {
	song := &domain.Song{
		ID:    id,
		Title: title,
		Genre: genre,
	}
	return u.songRepo.Update(ctx, song)
}

func (u *songUsecase) DeleteSong(ctx context.Context, id int64) error {
	return u.songRepo.Delete(ctx, id)
}

func (u *songUsecase) GetSongsByGenre(ctx context.Context, genre string, limit int32) ([]domain.Song, error) {
	return u.songRepo.GetByGenre(ctx, genre, limit)
}
