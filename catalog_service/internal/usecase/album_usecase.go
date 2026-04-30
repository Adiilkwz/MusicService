package usecase

import (
	"catalog_service/internal/domain"
	"context"
)

type albumUsecase struct {
	albumRepo domain.AlbumRepository
	songRepo  domain.SongRepository
}

func NewAlbumUsecase(al domain.AlbumRepository, sr domain.SongRepository) domain.AlbumUsecase {
	return &albumUsecase{
		albumRepo: al,
		songRepo:  sr,
	}
}

func (u *albumUsecase) CreateAlbum(ctx context.Context, artistID int64, title string, releaseYear int32) (int64, error) {
	album := &domain.Album{
		ArtistID:    artistID,
		Title:       title,
		ReleaseYear: releaseYear,
	}
	return u.albumRepo.Create(ctx, album)
}

func (u *albumUsecase) GetAlbum(ctx context.Context, id int64) (*domain.Album, error) {
	return u.albumRepo.GetByID(ctx, id)
}

func (u *albumUsecase) GetSongsByAlbum(ctx context.Context, albumID int64) ([]domain.Song, error) {
	return u.songRepo.GetByAlbumID(ctx, albumID)
}
