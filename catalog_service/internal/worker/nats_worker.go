package worker

import (
	"context"
	"encoding/json"
	"log"

	"catalog_service/internal/domain"

	"github.com/nats-io/nats.go"
)

type SongPlayedEvent struct {
	SongID int64 `json:"song_id"`
}

type NatsWorker struct {
	nc          *nats.Conn
	songUseCase domain.SongUsecase
}

func NewNatsWorker(nc *nats.Conn, sUC domain.SongUsecase) *NatsWorker {
	return &NatsWorker{
		nc:          nc,
		songUseCase: sUC,
	}
}

func (w *NatsWorker) Start(ctx context.Context) {
	_, err := w.nc.Subscribe("songs.played", func(m *nats.Msg) {
		var event SongPlayedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("NATS: ошибка парсинга сообщения: %v", err)
			return
		}

		log.Printf("NATS: получено событие прослушивания песни ID: %d", event.SongID)

		err := w.songUseCase.IncrementSongPlays(ctx, event.SongID)
		if err != nil {
			log.Printf("NATS: не удалось обновить счетчик для песни %d: %v", event.SongID, err)
		}
	})
	if err != nil {
		log.Fatalf("NATS: ошибка подписки: %v", err)
	}

	log.Println("NATS: воркер успешно запущен и слушает топик 'song.played'")
}
