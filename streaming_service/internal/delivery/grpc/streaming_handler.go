package grpc

import (
	"context"
	"io"
	"log"
	"strconv"

	"streaming_service/internal/metrics"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func extractUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "метаданные не найдены")
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return 0, status.Error(codes.Unauthenticated, "user_id отсутствует в метаданных")
	}

	userID, err := strconv.ParseInt(userIDs[0], 10, 64)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "неверный формат user_id")
	}

	return userID, nil
}

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
	metrics.GRPCRequests.WithLabelValues("RecordPlay").Inc()
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	err1 := h.streamingUC.RecordPlay(ctx, userID, req.SongId)
	if err1 != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to record play: %v", err1)
	}
	metrics.PlaysTotal.Inc()
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetUserHistory(ctx context.Context, req *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	history, err1 := h.streamingUC.GetUserHistory(ctx, userID, int(req.Limit))
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history: %v", err1)
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
