package domain

import "time"

// History represents a play history record
type History struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	SongID    int64     `json:"song_id"`
	PlayedAt  time.Time `json:"played_at"`
}

// Playlist represents a playlist
type Playlist struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Title  string  `json:"title"`
	SongIDs []int64 `json:"song_ids"`
}

// Like represents a song like
type Like struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	SongID int64 `json:"song_id"`
}

// TrendingItem represents a trending song
type TrendingItem struct {
	SongID    int64 `json:"song_id"`
	PlayCount int32 `json:"play_count"`
}