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
	cacheRepo  domain.CacheRepository
}

func NewSearchUsecase(ar domain.ArtistRepository, al domain.AlbumRepository, sr domain.SongRepository, cr domain.CacheRepository) domain.SearchUsecase {
	return &searchUsecase{
		artistRepo: ar,
		albumRepo:  al,
		songRepo:   sr,
		cacheRepo:  cr,
	}
}

func (u *searchUsecase) SearchCatalog(ctx context.Context, query string, limit int32) (*domain.SearchResult, error) {
	cacheKey := fmt.Sprintf("search:%s:limit:%d", query, limit)
	var result domain.SearchResult

	err := u.cacheRepo.Get(ctx, cacheKey, &result)
	if err == nil {
		log.Printf("ВНИМАНИЕ: Результат для '%s' взят из Redis (Кэш)!", query)
		return &result, nil
	}

	log.Printf("В кэше пусто. Ищем '%s' в PostgreSQL...", query)

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

	result = domain.SearchResult{
		Artists: artists,
		Albums:  albums,
		Songs:   songs,
	}

	if len(artists) > 0 || len(albums) > 0 || len(songs) > 0 {
		err = u.cacheRepo.Set(ctx, cacheKey, result, 300)
		if err != nil {
			log.Printf("Ошибка сохранения в кэш: %v", err) // Логируем, но не прерываем работу
		} else {
			log.Printf("Результат для '%s' успешно сохранен в Redis!", query)
		}
	}

	return &result, nil
}
