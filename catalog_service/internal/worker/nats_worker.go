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
	sub, err := w.nc.Subscribe("songs.played", func(m *nats.Msg) {
		var event SongPlayedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("NATS: ошибка парсинга сообщения: %v", err)
			return
		}

		log.Printf("NATS: получено событие прослушивания песни ID: %d", event.SongID)

		go func(songID int64) {
			if err := w.songUseCase.IncrementSongPlays(ctx, songID); err != nil {
				log.Printf("NATS: не удалось обновить счетчик для песни %d: %v", songID, err)
			}
		}(event.SongID)
	})
	if err != nil {
		log.Fatalf("NATS: ошибка подписки: %v", err)
	}

	log.Println("NATS: воркер успешно запущен и слушает топик 'songs.played'")

	go func() {
		<-ctx.Done()
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("NATS: ошибка при отписке: %v", err)
		} else {
			log.Println("NATS: отписан и завершаю работу воркера")
		}
	}()
}
