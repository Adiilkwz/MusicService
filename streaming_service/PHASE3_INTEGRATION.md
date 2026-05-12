# ФАЗА 3: Внутренняя коммуникация

## Что реализовано в Streaming Service (Participant 3)

### 1. Event Publisher - NATS для публикации событий
- ✅ `streaming_service/internal/events/events.go` - структура `SongPlayedEvent`
- ✅ `streaming_service/internal/repository/nats/event_publisher.go` - публикатор событий
- ✅ При `RecordPlay` отправляется асинхронное событие в NATS (subject: `songs.played`)
- ✅ Событие содержит: `user_id`, `song_id`, `timestamp`, `event_id`

### 2. Event Subscriber (ПРИМЕР для других сервисов)
- ✅ `streaming_service/internal/repository/nats/event_subscriber.go` - пример подписчика

---

## ЧТО НУЖНО СДЕЛАТЬ ДРУГИМ УЧАСТНИКАМ

### Participant 1 (Auth Service)
```go
// internal/delivery/grpc/validate_handler.go
func (s *authServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    // Проверить JWT токен
    // Вернуть user_id, roles и expiry
    // Это будет "паспорт" для других сервисов
}
```

**Endpoint:** `auth_service:50050` → `AuthService.ValidateToken`

---

### Participant 2 (Catalog Service)
```go
// internal/repository/nats/song_event_handler.go
func (h *SongEventHandler) HandleSongPlayed(event *events.SongPlayedEvent) error {
    // Инкрементировать play_count для песни в БД или Redis
    // UPDATE songs SET play_count = play_count + 1 WHERE id = event.SongId
    // Или в Redis: INCR songs:{song_id}:play_count
}

// main.go инициализация
subscriber, _ := nats.NewEventSubscriber(cfg.NATSUrl)
subscriber.OnSongPlayed(eventHandler.HandleSongPlayed)
```

---

## Пример интеграции

```bash
# Streaming Service публикует
Event: SongPlayedEvent {
  user_id: 123,
  song_id: 456,
  timestamp: "2026-05-10T10:30:00Z"
}

# Catalog Service получает и обновляет
UPDATE songs SET play_count = play_count + 1 WHERE id = 456
```

---

## Команды для интеграции

```bash
# Streaming Service - опубликовать событие
# (автоматически при RecordPlay)

# Catalog Service - подписаться
nc = natsconnect(url)
nc.subscribe("songs.played", handler)

# Auth Service - запустить ValidateToken gRPC endpoint
grpc_server.Register(AuthService)
```