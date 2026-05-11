package usecase

import (
	"context"

	"streaming_service/internal/domain"
)

type playlistUsecase struct {
	playlistRepo domain.PlaylistRepository
}

func NewPlaylistUsecase(playlistRepo domain.PlaylistRepository) domain.PlaylistUsecase {
	return &playlistUsecase{playlistRepo: playlistRepo}
}

func (u *playlistUsecase) CreatePlaylist(ctx context.Context, userID int64, title string) (int64, error) {
	if title == "" {
		title = "Untitled Playlist"
	}
	return u.playlistRepo.Create(ctx, userID, title)
}

func (u *playlistUsecase) GetPlaylist(ctx context.Context, playlistID int64) (*domain.Playlist, error) {
	return u.playlistRepo.GetByID(ctx, playlistID)
}

func (u *playlistUsecase) AddSongToPlaylist(ctx context.Context, playlistID, songID int64) error {
	return u.playlistRepo.AddSong(ctx, playlistID, songID)
}

func (u *playlistUsecase) RemoveSongFromPlaylist(ctx context.Context, playlistID, songID int64) error {
	return u.playlistRepo.RemoveSong(ctx, playlistID, songID)
}

func (u *playlistUsecase) DeletePlaylist(ctx context.Context, playlistID int64) error {
	return u.playlistRepo.Delete(ctx, playlistID)
}
