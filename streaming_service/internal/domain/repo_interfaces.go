package domain

import "context"

type HistoryRepository interface {
	Create(ctx context.Context, userID, songID int64) error
	GetByUserID(ctx context.Context, userID int64, limit int) ([]History, error)
}

type PlaylistRepository interface {
	Create(ctx context.Context, userID int64, title string) (int64, error)
	GetByID(ctx context.Context, playlistID int64) (*Playlist, error)
	AddSong(ctx context.Context, playlistID, songID int64) error
	RemoveSong(ctx context.Context, playlistID, songID int64) error
	Delete(ctx context.Context, playlistID int64) error
}

type LikeRepository interface {
	Like(ctx context.Context, userID, songID int64) error
	Unlike(ctx context.Context, userID, songID int64) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]int64, error)
}

type TrendingRepository interface {
	GetTrending(ctx context.Context, limit int) ([]TrendingItem, error)
	IncrementPlayCount(ctx context.Context, songID int64) error
}

type AudioRepository interface {
	GetAudioPath(songID int64) string
	ReadChunk(songID int64, offset, size int64) ([]byte, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expirationSeconds int) error
	Get(ctx context.Context, key string, dest interface{}) error
}
