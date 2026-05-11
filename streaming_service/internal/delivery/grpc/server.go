package grpc

import (
	"streaming_service/internal/domain"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc"
)

type Handler struct {
	pb.UnimplementedStreamingServiceServer
	streamingUC domain.StreamingUsecase
	playlistUC  domain.PlaylistUsecase
	likeUC      domain.LikeUsecase
}

func NewHandler(
	streamingUC domain.StreamingUsecase,
	playlistUC domain.PlaylistUsecase,
	likeUC domain.LikeUsecase,
) *Handler {
	return &Handler{
		streamingUC: streamingUC,
		playlistUC:  playlistUC,
		likeUC:      likeUC,
	}
}

func Register(server *grpc.Server, h *Handler) {
	pb.RegisterStreamingServiceServer(server, h)
}
