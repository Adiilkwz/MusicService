package events

import (
	"encoding/json"
	"time"
)

// SongPlayedEvent represents when a user finishes playing a song
type SongPlayedEvent struct {
	EventID   string    `json:"event_id"`
	UserID    int64     `json:"user_id"`
	SongID    int64     `json:"song_id"`
	Timestamp time.Time `json:"timestamp"`
	PlayCount int32     `json:"play_count"` // Optional: incremented count
}

// Marshal serializes event to JSON
func (e *SongPlayedEvent) Marshal() []byte {
	data, _ := json.Marshal(e)
	return data
}

// Subject returns NATS subject for this event
func (e *SongPlayedEvent) Subject() string {
	return "songs.played"
}