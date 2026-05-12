package usecase

import (
	"context"
	"fmt"
	"log"

	"catalog_service/internal/domain"
)

type searchUsecase struct {
	artistRepo domain.ArtistRepository
	albumRepo  domain.AlbumRepository
	songRepo   domain.SongRepository
}

func NewSearchUsecase(ar domain.ArtistRepository, al domain.AlbumRepository, sr domain.SongRepository) domain.SearchUsecase {
	return &searchUsecase{
		artistRepo: ar,
		albumRepo:  al,
		songRepo:   sr,
	}
}

func (u *searchUsecase) SearchCatalog(ctx context.Context, query string, limit int32) (*domain.SearchResult, error) {
	log.Printf("Ищем '%s' напрямую в PostgreSQL...", query)

	artists, err := u.artistRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска артистов: %w", err)
	}

	albums, err := u.albumRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска альбомов: %w", err)
	}

	songs, err := u.songRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска песен: %w", err)
	}

	return &domain.SearchResult{
		Artists: artists,
		Albums:  albums,
		Songs:   songs,
	}, nil
}
