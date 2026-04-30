package usecase

import (
	"context"
	"io"
	"os"

	"streaming_service/internal/domain"
)

type streamingUsecase struct {
	historyRepo   domain.HistoryRepository
	trendingRepo   domain.TrendingRepository
	audioRepo      domain.AudioRepository
}

func NewStreamingUsecase(
	historyRepo domain.HistoryRepository,
	trendingRepo domain.TrendingRepository,
	audioRepo domain.AudioRepository,
) domain.StreamingUsecase {
	return &streamingUsecase{
		historyRepo:  historyRepo,
		trendingRepo: trendingRepo,
		audioRepo:    audioRepo,
	}
}

func (u *streamingUsecase) StreamAudio(ctx context.Context, songID int64) (io.ReadCloser, error) {
	audioPath := u.audioRepo.GetAudioPath(songID)
	
	// Check if file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		// Return empty reader if file not found (for demo purposes)
		return &emptyReader{}, nil
	}
	
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, err
	}
	
	return file, nil
}

func (u *streamingUsecase) RecordPlay(ctx context.Context, userID, songID int64) error {
	// Record play in history
	if err := u.historyRepo.Create(ctx, userID, songID); err != nil {
		return err
	}
	
	// Increment trending play count
	return u.trendingRepo.IncrementPlayCount(ctx, songID)
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

type playlistUsecase struct {
	playlistRepo domain.PlaylistRepository
}

func NewPlaylistUsecase(playlistRepo domain.PlaylistRepository) domain.PlaylistUsecase {
	return &playlistUsecase{playlistRepo: playlistRepo}
}

func (u *playlistUsecase) CreatePlaylist(ctx context.Context, userID int64, title string) (int64, error) {
	if title == "" {
		title = "Untitled Playlist"
	}
	return u.playlistRepo.Create(ctx, userID, title)
}

func (u *playlistUsecase) GetPlaylist(ctx context.Context, playlistID int64) (*domain.Playlist, error) {
	return u.playlistRepo.GetByID(ctx, playlistID)
}

func (u *playlistUsecase) AddSongToPlaylist(ctx context.Context, playlistID, songID int64) error {
	return u.playlistRepo.AddSong(ctx, playlistID, songID)
}

func (u *playlistUsecase) RemoveSongFromPlaylist(ctx context.Context, playlistID, songID int64) error {
	return u.playlistRepo.RemoveSong(ctx, playlistID, songID)
}

func (u *playlistUsecase) DeletePlaylist(ctx context.Context, playlistID int64) error {
	return u.playlistRepo.Delete(ctx, playlistID)
}

type likeUsecase struct {
	likeRepo domain.LikeRepository
}

func NewLikeUsecase(likeRepo domain.LikeRepository) domain.LikeUsecase {
	return &likeUsecase{likeRepo: likeRepo}
}

func (u *likeUsecase) LikeSong(ctx context.Context, userID, songID int64) error {
	return u.likeRepo.Like(ctx, userID, songID)
}

func (u *likeUsecase) UnlikeSong(ctx context.Context, userID, songID int64) error {
	return u.likeRepo.Unlike(ctx, userID, songID)
}

func (u *likeUsecase) GetLikedSongs(ctx context.Context, userID int64, limit, offset int) ([]int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return u.likeRepo.GetByUserID(ctx, userID, limit, offset)
}