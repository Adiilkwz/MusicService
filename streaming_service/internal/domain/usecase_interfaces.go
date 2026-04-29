package domain

import (
	"context"
	"io"
)

// StreamingUsecase defines streaming business logic
type StreamingUsecase interface {
	StreamAudio(ctx context.Context, songID int64) (io.ReadCloser, error)
	RecordPlay(ctx context.Context, userID, songID int64) error
	GetUserHistory(ctx context.Context, userID int64, limit int) ([]History, error)
	GetTrending(ctx context.Context, limit int) ([]TrendingItem, error)
}

// PlaylistUsecase defines playlist business logic
type PlaylistUsecase interface {
	CreatePlaylist(ctx context.Context, userID int64, title string) (int64, error)
	GetPlaylist(ctx context.Context, playlistID int64) (*Playlist, error)
	AddSongToPlaylist(ctx context.Context, playlistID, songID int64) error
	RemoveSongFromPlaylist(ctx context.Context, playlistID, songID int64) error
	DeletePlaylist(ctx context.Context, playlistID int64) error
}

// LikeUsecase defines like business logic
type LikeUsecase interface {
	LikeSong(ctx context.Context, userID, songID int64) error
	UnlikeSong(ctx context.Context, userID, songID int64) error
	GetLikedSongs(ctx context.Context, userID int64, limit, offset int) ([]int64, error)
}