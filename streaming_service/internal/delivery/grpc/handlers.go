package grpc

import (
	"context"
	"io"
	"log"

	"streaming_service/internal/domain"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamingServer struct {
	pb.UnimplementedStreamingServiceServer
	usecase domain.StreamingUsecase
}

func NewStreamingServer(usecase domain.StreamingUsecase) *StreamingServer {
	return &StreamingServer{usecase: usecase}
}

func (s *StreamingServer) StreamAudio(req *pb.StreamRequest, stream pb.StreamingService_StreamAudioServer) error {
	songID := req.SongId

	reader, err := s.usecase.StreamAudio(stream.Context(), songID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to stream audio: %v", err)
	}
	defer reader.Close()

	buf := make([]byte, 4096) // 4KB chunks
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			response := &pb.StreamResponse{
				AudioChunk: buf[:n],
			}
			if err := stream.Send(response); err != nil {
				log.Printf("Error sending chunk: %v", err)
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "error reading audio: %v", err)
		}
	}

	return nil
}

func (s *StreamingServer) RecordPlay(ctx context.Context, req *pb.RecordPlayRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.RecordPlay(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to record play: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *StreamingServer) GetUserHistory(ctx context.Context, req *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	history, err := s.usecase.GetUserHistory(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history: %v", err)
	}

	var items []*pb.HistoryItem
	for _, h := range history {
		items = append(items, &pb.HistoryItem{
			SongId:   h.SongID,
			PlayedAt: h.PlayedAt.String(),
		})
	}

	return &pb.GetUserHistoryResponse{History: items}, nil
}

func (s *StreamingServer) GetTrending(ctx context.Context, req *pb.GetTrendingRequest) (*pb.GetTrendingResponse, error) {
	trending, err := s.usecase.GetTrending(ctx, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get trending: %v", err)
	}

	var items []*pb.TrendingItem
	for _, t := range trending {
		items = append(items, &pb.TrendingItem{
			SongId:    t.SongID,
			PlayCount: t.PlayCount,
		})
	}

	return &pb.GetTrendingResponse{Items: items}, nil
}

// PlaylistServiceServer implementation
type PlaylistServer struct {
	pb.UnimplementedStreamingServiceServer
	usecase domain.PlaylistUsecase
}

func NewPlaylistServer(usecase domain.PlaylistUsecase) *PlaylistServer {
	return &PlaylistServer{usecase: usecase}
}

func (s *PlaylistServer) CreatePlaylist(ctx context.Context, req *pb.CreatePlaylistRequest) (*pb.CreatePlaylistResponse, error) {
	playlistID, err := s.usecase.CreatePlaylist(ctx, req.UserId, req.Title)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create playlist: %v", err)
	}
	return &pb.CreatePlaylistResponse{PlaylistId: playlistID}, nil
}

func (s *PlaylistServer) GetPlaylist(ctx context.Context, req *pb.GetPlaylistRequest) (*pb.GetPlaylistResponse, error) {
	playlist, err := s.usecase.GetPlaylist(ctx, req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get playlist: %v", err)
	}
	return &pb.GetPlaylistResponse{
		PlaylistId: playlist.ID,
		Title:      playlist.Title,
		SongIds:    playlist.SongIDs,
	}, nil
}

func (s *PlaylistServer) AddSongToPlaylist(ctx context.Context, req *pb.ModifyPlaylistRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.AddSongToPlaylist(ctx, req.PlaylistId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to add song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *PlaylistServer) RemoveSongFromPlaylist(ctx context.Context, req *pb.ModifyPlaylistRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.RemoveSongFromPlaylist(ctx, req.PlaylistId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to remove song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *PlaylistServer) DeletePlaylist(ctx context.Context, req *pb.DeletePlaylistRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.DeletePlaylist(ctx, req.PlaylistId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to delete playlist: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

// LikeServiceServer implementation
type LikeServer struct {
	pb.UnimplementedStreamingServiceServer
	usecase domain.LikeUsecase
}

func NewLikeServer(usecase domain.LikeUsecase) *LikeServer {
	return &LikeServer{usecase: usecase}
}

func (s *LikeServer) LikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.LikeSong(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to like song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *LikeServer) UnlikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	err := s.usecase.UnlikeSong(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to unlike song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *LikeServer) GetLikedSongs(ctx context.Context, req *pb.GetLikedSongsRequest) (*pb.GetLikedSongsResponse, error) {
	songIDs, err := s.usecase.GetLikedSongs(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked songs: %v", err)
	}
	return &pb.GetLikedSongsResponse{SongIds: songIDs}, nil
}
