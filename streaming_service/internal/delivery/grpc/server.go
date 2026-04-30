package grpc

import (
	"streaming_service/internal/domain"
	"google.golang.org/grpc"
)

type Server struct {
	streamingServer *StreamingServer
	playlistServer   *PlaylistServer
	likeServer       *LikeServer
}

func NewServer(
	streamingUsecase domain.StreamingUsecase,
	playlistUsecase domain.PlaylistUsecase,
	likeUsecase domain.LikeUsecase,
) *Server {
	return &Server{
		streamingServer: NewStreamingServer(streamingUsecase),
		playlistServer:   NewPlaylistServer(playlistUsecase),
		likeServer:       NewLikeServer(likeUsecase),
	}
}

func (s *Server) Register(server *grpc.Server) {
	RegisterStreamingServiceServer(server, s.streamingServer)
	RegisterPlaylistServiceServer(server, s.playlistServer)
	RegisterLikeServiceServer(server, s.likeServer)
}