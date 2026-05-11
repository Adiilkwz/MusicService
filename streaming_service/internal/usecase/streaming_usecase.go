package usecase

import (
	"context"
	"io"
	"os"
	"time"

	"streaming_service/internal/domain"
)

type streamingUsecase struct {
	historyRepo    domain.HistoryRepository
	trendingRepo   domain.TrendingRepository
	audioRepo      domain.AudioRepository
	eventPublisher domain.EventPublisher
}

func NewStreamingUsecase(
	historyRepo domain.HistoryRepository,
	trendingRepo domain.TrendingRepository,
	audioRepo domain.AudioRepository,
	eventPublisher domain.EventPublisher,
) domain.StreamingUsecase {
	return &streamingUsecase{
		historyRepo:    historyRepo,
		trendingRepo:   trendingRepo,
		audioRepo:      audioRepo,
		eventPublisher: eventPublisher,
	}
}

func (u *streamingUsecase) StreamAudio(ctx context.Context, songID int64) (io.ReadCloser, error) {
	audioPath := u.audioRepo.GetAudioPath(songID)

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return &emptyReader{}, nil
	}

	file, err := os.Open(audioPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (u *streamingUsecase) RecordPlay(ctx context.Context, userID, songID int64) error {
	if err := u.historyRepo.Create(ctx, userID, songID); err != nil {
		return err
	}

	if err := u.trendingRepo.IncrementPlayCount(ctx, songID); err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := u.eventPublisher.PublishSongPlayed(ctx, userID, songID); err != nil {
			return
		}
	}()

	return nil
}

func (u *streamingUsecase) GetUserHistory(ctx context.Context, userID int64, limit int) ([]domain.History, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.historyRepo.GetByUserID(ctx, userID, limit)
}

func (u *streamingUsecase) GetTrending(ctx context.Context, limit int) ([]domain.TrendingItem, error) {
	if limit <= 0 {
		limit = 10
	}
	return u.trendingRepo.GetTrending(ctx, limit)
}

type emptyReader struct{}

func (r *emptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (r *emptyReader) Close() error {
	return nil
}
