package usecase

import (
	"catalog_service/internal/domain"
	"context"
)

type artistUsecase struct {
	artistRepo domain.ArtistRepository
	albumRepo  domain.AlbumRepository
}

func NewArtistUsecase(ar domain.ArtistRepository, al domain.AlbumRepository) domain.ArtistUsecase {
	return &artistUsecase{
		artistRepo: ar,
		albumRepo:  al,
	}
}

func (u *artistUsecase) CreateArtist(ctx context.Context, name, bio string) (int64, error) {
	artist := &domain.Artist{
		Name: name,
		Bio:  bio,
	}
	return u.artistRepo.Create(ctx, artist)
}

func (u *artistUsecase) GetArtist(ctx context.Context, id int64) (*domain.Artist, error) {
	return u.artistRepo.GetByID(ctx, id)
}

func (u *artistUsecase) GetAlbumsByArtist(ctx context.Context, artistID int64) ([]domain.Album, error) {
	return u.albumRepo.GetByArtistID(ctx, artistID)
}
