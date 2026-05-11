package grpc

import (
	"context"
	"io"
	"log"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) StreamAudio(req *pb.StreamRequest, stream pb.StreamingService_StreamAudioServer) error {
	reader, err := h.streamingUC.StreamAudio(stream.Context(), req.SongId)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to stream audio: %v", err)
	}
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			response := &pb.StreamResponse{AudioChunk: buf[:n]}
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

func (h *Handler) RecordPlay(ctx context.Context, req *pb.RecordPlayRequest) (*pb.SuccessResponse, error) {
	err := h.streamingUC.RecordPlay(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to record play: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetUserHistory(ctx context.Context, req *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	history, err := h.streamingUC.GetUserHistory(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history: %v", err)
	}

	var items []*pb.HistoryItem
	for _, hist := range history {
		items = append(items, &pb.HistoryItem{
			SongId:   hist.SongID,
			PlayedAt: hist.PlayedAt.String(),
		})
	}
	return &pb.GetUserHistoryResponse{History: items}, nil
}

func (h *Handler) GetTrending(ctx context.Context, req *pb.GetTrendingRequest) (*pb.GetTrendingResponse, error) {
	trending, err := h.streamingUC.GetTrending(ctx, int(req.Limit))
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
