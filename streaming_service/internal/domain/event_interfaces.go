package domain

import "context"

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	PublishSongPlayed(ctx context.Context, userID, songID int64) error
	Close() error
}