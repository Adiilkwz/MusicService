package grpc

import (
	"context"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) LikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	err := h.likeUC.LikeSong(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to like song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) UnlikeSong(ctx context.Context, req *pb.LikeSongRequest) (*pb.SuccessResponse, error) {
	err := h.likeUC.UnlikeSong(ctx, req.UserId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to unlike song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetLikedSongs(ctx context.Context, req *pb.GetLikedSongsRequest) (*pb.GetLikedSongsResponse, error) {
	songIDs, err := h.likeUC.GetLikedSongs(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked songs: %v", err)
	}
	return &pb.GetLikedSongsResponse{SongIds: songIDs}, nil
}
