																																																																																																																															package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"streaming_service/config"
	"streaming_service/internal/repository/postgres"
	"streaming_service/internal/usecase"
)

// Mock types for proto messages
type SuccessResponse struct {
	Success bool
}

type StreamRequest struct {
	SongId int64
}

type StreamResponse struct {
	AudioChunk []byte
}

type RecordPlayRequest struct {
	UserId int64
	SongId int64
}

type GetUserHistoryRequest struct {
	UserId int64
	Limit  int32
}

type HistoryItem struct {
	SongId   int64
	PlayedAt string
}

type GetUserHistoryResponse struct {
	History []*HistoryItem
}

type GetTrendingRequest struct {
	Limit int32
}

type TrendingItem struct {
	SongId    int64
	PlayCount int32
}

type GetTrendingResponse struct {
	Items []*TrendingItem
}

type CreatePlaylistRequest struct {
	UserId int64
	Title  string
}

type CreatePlaylistResponse struct {
	PlaylistId int64
}

type GetPlaylistRequest struct {
	PlaylistId int64
}

type GetPlaylistResponse struct {
	PlaylistId int64
	Title      string
	SongIds    []int64
}

type ModifyPlaylistRequest struct {
	PlaylistId int64
	SongId     int64
}

type DeletePlaylistRequest struct {
	PlaylistId int64
}

type LikeSongRequest struct {
	UserId int64
	SongId int64
}

type GetLikedSongsRequest struct {
	UserId int64
	Limit  int32
	Offset int32
}

type GetLikedSongsResponse struct {
	SongIds []int64
}

// StreamingServiceServer interface
type StreamingServiceServer interface {
	StreamAudio(*StreamRequest, StreamingService_StreamAudioServer) error
	RecordPlay(context.Context, *RecordPlayRequest) (*SuccessResponse, error)
	GetUserHistory(context.Context, *GetUserHistoryRequest) (*GetUserHistoryResponse, error)
	GetTrending(context.Context, *GetTrendingRequest) (*GetTrendingResponse, error)
}

type StreamingService_StreamAudioServer interface {
	Send(*StreamResponse) error
	Context() context.Context
}

// PlaylistServiceServer interface
type PlaylistServiceServer interface {
	CreatePlaylist(context.Context, *CreatePlaylistRequest) (*CreatePlaylistResponse, error)
	GetPlaylist(context.Context, *GetPlaylistRequest) (*GetPlaylistResponse, error)
	AddSongToPlaylist(context.Context, *ModifyPlaylistRequest) (*SuccessResponse, error)
	RemoveSongFromPlaylist(context.Context, *ModifyPlaylistRequest) (*SuccessResponse, error)
	DeletePlaylist(context.Context, *DeletePlaylistRequest) (*SuccessResponse, error)
}

// LikeServiceServer interface
type LikeServiceServer interface {
	LikeSong(context.Context, *LikeSongRequest) (*SuccessResponse, error)
	UnlikeSong(context.Context, *LikeSongRequest) (*SuccessResponse, error)
	GetLikedSongs(context.Context, *GetLikedSongsRequest) (*GetLikedSongsResponse, error)
}

// StreamingServer implementation
type streamingServer struct {
	usecase usecase.StreamingUsecase
}

func (s *streamingServer) StreamAudio(req *StreamRequest, stream StreamingService_StreamAudioServer) error {
	return status.Errorf(codes.Unimplemented, "streaming not implemented")
}

func (s *streamingServer) RecordPlay(ctx context.Context, req *RecordPlayRequest) (*SuccessResponse, error) {
	err := s.usecase.RecordPlay(ctx, req.UserId, req.SongId)
	return &SuccessResponse{Success: err == nil}, err
}

func (s *streamingServer) GetUserHistory(ctx context.Context, req *GetUserHistoryRequest) (*GetUserHistoryResponse, error) {
	history, err := s.usecase.GetUserHistory(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, err
	}
	var items []*HistoryItem
	for _, h := range history {
		items = append(items, &HistoryItem{
			SongId:   h.SongID,
			PlayedAt: h.PlayedAt.String(),
		})
	}
	return &GetUserHistoryResponse{History: items}, nil
}

func (s *streamingServer) GetTrending(ctx context.Context, req *GetTrendingRequest) (*GetTrendingResponse, error) {
	trending, err := s.usecase.GetTrending(ctx, int(req.Limit))
	if err != nil {
		return nil, err
	}
	var items []*TrendingItem
	for _, t := range trending {
		items = append(items, &TrendingItem{
			SongId:    t.SongID,
			PlayCount: t.PlayCount,
		})
	}
	return &GetTrendingResponse{Items: items}, nil
}

// PlaylistServer implementation
type playlistServer struct {
	usecase usecase.PlaylistUsecase
}

func (s *playlistServer) CreatePlaylist(ctx context.Context, req *CreatePlaylistRequest) (*CreatePlaylistResponse, error) {
	id, err := s.usecase.CreatePlaylist(ctx, req.UserId, req.Title)
	return &CreatePlaylistResponse{PlaylistId: id}, err
}

func (s *playlistServer) GetPlaylist(ctx context.Context, req *GetPlaylistRequest) (*GetPlaylistResponse, error) {
	p, err := s.usecase.GetPlaylist(ctx, req.PlaylistId)
	if err != nil {
		return nil, err
	}
	return &GetPlaylistResponse{PlaylistId: p.ID, Title: p.Title, SongIds: p.SongIDs}, nil
}

func (s *playlistServer) AddSongToPlaylist(ctx context.Context, req *ModifyPlaylistRequest) (*SuccessResponse, error) {
	err := s.usecase.AddSongToPlaylist(ctx, req.PlaylistId, req.SongId)
	return &SuccessResponse{Success: err == nil}, err
}

func (s *playlistServer) RemoveSongFromPlaylist(ctx context.Context, req *ModifyPlaylistRequest) (*SuccessResponse, error) {
	err := s.usecase.RemoveSongFromPlaylist(ctx, req.PlaylistId, req.SongId)
	return &SuccessResponse{Success: err == nil}, err
}

func (s *playlistServer) DeletePlaylist(ctx context.Context, req *DeletePlaylistRequest) (*SuccessResponse, error) {
	err := s.usecase.DeletePlaylist(ctx, req.PlaylistId)
	return &SuccessResponse{Success: err == nil}, err
}

// LikeServer implementation
type likeServer struct {
	usecase usecase.LikeUsecase
}

func (s *likeServer) LikeSong(ctx context.Context, req *LikeSongRequest) (*SuccessResponse, error) {
	err := s.usecase.LikeSong(ctx, req.UserId, req.SongId)
	return &SuccessResponse{Success: err == nil}, err
}

func (s *likeServer) UnlikeSong(ctx context.Context, req *LikeSongRequest) (*SuccessResponse, error) {
	err := s.usecase.UnlikeSong(ctx, req.UserId, req.SongId)
	return &SuccessResponse{Success: err == nil}, err
}

func (s *likeServer) GetLikedSongs(ctx context.Context, req *GetLikedSongsRequest) (*GetLikedSongsResponse, error) {
	songIDs, err := s.usecase.GetLikedSongs(ctx, req.UserId, int(req.Limit), int(req.Offset))
	return &GetLikedSongsResponse{SongIds: songIDs}, err
}

func main() {
	cfg := config.Load()

	db, err := postgres.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	historyRepo := postgres.NewHistoryRepository(db)
	playlistRepo := postgres.NewPlaylistRepository(db)
	likeRepo := postgres.NewLikeRepository(db)
	trendingRepo := postgres.NewTrendingRepository(db)
	audioRepo := postgres.NewAudioRepository(cfg.AudioDir)

	streamingUsecase := usecase.NewStreamingUsecase(historyRepo, trendingRepo, audioRepo)
	playlistUsecase := usecase.NewPlaylistUsecase(playlistRepo)
	likeUsecase := usecase.NewLikeUsecase(likeRepo)

	// Register handlers
	streamingHandler := &streamingServer{usecase: streamingUsecase}
	playlistHandler := &playlistServer{usecase: playlistUsecase}
	likeHandler := &likeServer{usecase: likeUsecase}

	server := grpc.NewServer()

	// Note: In real implementation, these would be registered via pb.RegisterXXX
	// For now, handlers are ready to be registered when proto is generated

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Streaming service starting on port %s", cfg.ServerPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Use handlers to avoid unused variable errors
	_ = streamingHandler
	_ = playlistHandler
	_ = likeHandler
}