package domain

import "context"

// HistoryRepository defines the interface for history data access
type HistoryRepository interface {
	Create(ctx context.Context, userID, songID int64) error
	GetByUserID(ctx context.Context, userID int64, limit int) ([]History, error)
}

// PlaylistRepository defines the interface for playlist data access
type PlaylistRepository interface {
	Create(ctx context.Context, userID int64, title string) (int64, error)
	GetByID(ctx context.Context, playlistID int64) (*Playlist, error)
	AddSong(ctx context.Context, playlistID, songID int64) error
	RemoveSong(ctx context.Context, playlistID, songID int64) error
	Delete(ctx context.Context, playlistID int64) error
}

// LikeRepository defines the interface for likes data access
type LikeRepository interface {
	Like(ctx context.Context, userID, songID int64) error
	Unlike(ctx context.Context, userID, songID int64) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]int64, error)
}

// TrendingRepository defines the interface for trending data access
type TrendingRepository interface {
	GetTrending(ctx context.Context, limit int) ([]TrendingItem, error)
	IncrementPlayCount(ctx context.Context, songID int64) error
}

// AudioRepository defines the interface for audio file access
type AudioRepository interface {
	GetAudioPath(songID int64) string
	ReadChunk(songID int64, offset, size int64) ([]byte, error)
}