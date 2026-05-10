package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"streaming_service/internal/events"
)

// EventSubscriber listens to NATS events
// This is an EXAMPLE for other services (like Catalog service)
type EventSubscriber struct {
	conn *nats.Conn
	sub  *nats.Subscription
}

// NewEventSubscriber creates a new NATS subscriber
// EXAMPLE: Used by Catalog service to listen to SongPlayedEvent
func NewEventSubscriber(natsURL string) (*EventSubscriber, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &EventSubscriber{conn: conn}, nil
}

// OnSongPlayed subscribes to song played events
// EXAMPLE: Catalog increments play_count when receiving this event
func (s *EventSubscriber) OnSongPlayed(handler func(*events.SongPlayedEvent) error) error {
	sub, err := s.conn.Subscribe("songs.played", func(msg *nats.Msg) {
		var event events.SongPlayedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Error handling event: %v", err)
		}
	})

	if err != nil {
		return err
	}

	s.sub = sub
	return nil
}

// Close closes the subscriber
func (s *EventSubscriber) Close() error {
	if s.sub != nil {
		s.sub.Unsubscribe()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}