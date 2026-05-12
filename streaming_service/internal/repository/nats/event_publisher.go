package nats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"streaming_service/internal/events"
)

type eventPublisher struct {
	conn *nats.Conn
}

// NewEventPublisher creates a new NATS event publisher
func NewEventPublisher(natsURL string) (*eventPublisher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &eventPublisher{conn: conn}, nil
}

// PublishSongPlayed publishes a SongPlayedEvent to NATS
func (p *eventPublisher) PublishSongPlayed(ctx context.Context, userID, songID int64) error {
	event := &events.SongPlayedEvent{
		EventID:   fmt.Sprintf("%d-%d-%d", userID, songID, time.Now().UnixNano()),
		UserID:    userID,
		SongID:    songID,
		Timestamp: time.Now(),
	}

	data := event.Marshal()
	
	// Publish asynchronously with timeout
	done := make(chan error, 1)
	go func() {
		err := p.conn.Publish(event.Subject(), data)
		if err != nil {
			log.Printf("Error publishing SongPlayedEvent: %v", err)
		}
		done <- err
	}()

	// Wait for publish with context timeout
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("publish timeout")
	}
}

// Close closes the NATS connection
func (p *eventPublisher) Close() error {
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}