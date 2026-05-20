package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"auth_service/internal/domain"

	"github.com/nats-io/nats.go"
)

type natsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) domain.EventPublisher {
	return &natsPublisher{nc: nc}
}

type UserRegisteredEvent struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func (p *natsPublisher) PublishUserRegistered(userID int64, email, displayName string) error {
	event := UserRegisteredEvent{
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.nc.Publish("user.events.registered", data)
	if err != nil {
		return fmt.Errorf("failed to publish nats message: %w", err)
	}

	log.Printf("[NATS PUBLISHER] Broadcasted new user registration: %s\n", email)
	return nil
}
