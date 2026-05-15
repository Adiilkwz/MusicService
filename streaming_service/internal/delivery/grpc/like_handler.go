package grpc

import (
	"context"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) LikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	err1 := h.likeUC.LikeSong(ctx, userID, req.SongId)
	if err1 != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to like song: %v", err1)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) UnlikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	err1 := h.likeUC.UnlikeSong(ctx, userID, req.SongId)
	if err1 != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to unlike song: %v", err1)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetLikedSongs(ctx context.Context, req *pb.GetLikedSongsRequest) (*pb.GetLikedSongsResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	songIDs, err1 := h.likeUC.GetLikedSongs(ctx, userID, int(req.Limit), int(req.Offset))
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked songs: %v", err1)
	}
	return &pb.GetLikedSongsResponse{SongIds: songIDs}, nil
}
