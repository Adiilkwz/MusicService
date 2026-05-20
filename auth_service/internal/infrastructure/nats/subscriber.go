package nats

import (
	"encoding/json"
	"log"

	"auth_service/internal/domain"

	"github.com/nats-io/nats.go"
)

func StartRegistrationEmailWorker(nc *nats.Conn, emailSender domain.EmailSender) {
	_, err := nc.Subscribe("user.events.registered", func(m *nats.Msg) {
		log.Println("[NATS WORKER] Received registration event! Sending welcome email...")

		var event UserRegisteredEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("[NATS WORKER] Failed to decode message data: %v\n", err)
			return
		}

		err := emailSender.SendWelcomeEmail(event.Email, event.DisplayName)
		if err != nil {
			log.Printf("[NATS WORKER] Failed to send email to %s: %v\n", event.Email, err)
		} else {
			log.Printf("[NATS WORKER] Welcome email successfully sent to %s\n", event.Email)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS: %v", err)
	}

	log.Println("[NATS WORKER] Listening for registration events...")
}
