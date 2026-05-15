package grpc

import (
	"context"

	pb "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreatePlaylist(ctx context.Context, req *pb.CreatePlaylistRequest) (*pb.CreatePlaylistResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	playlistID, err1 := h.playlistUC.CreatePlaylist(ctx, userID, req.Title)
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "failed to create playlist: %v", err1)
	}
	return &pb.CreatePlaylistResponse{PlaylistId: playlistID}, nil
}

func (h *Handler) GetPlaylist(ctx context.Context, req *pb.GetPlaylistRequest) (*pb.GetPlaylistResponse, error) {
	playlist, err := h.playlistUC.GetPlaylist(ctx, req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get playlist: %v", err)
	}
	return &pb.GetPlaylistResponse{
		PlaylistId: playlist.ID,
		Title:      playlist.Title,
		SongIds:    playlist.SongIDs,
	}, nil
}

func (h *Handler) AddSongToPlaylist(ctx context.Context, req *pb.ModifyPlaylistRequest) (*pb.SuccessResponse, error) {
	err := h.playlistUC.AddSongToPlaylist(ctx, req.PlaylistId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to add song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) RemoveSongFromPlaylist(ctx context.Context, req *pb.ModifyPlaylistRequest) (*pb.SuccessResponse, error) {
	err := h.playlistUC.RemoveSongFromPlaylist(ctx, req.PlaylistId, req.SongId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to remove song: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (h *Handler) DeletePlaylist(ctx context.Context, req *pb.DeletePlaylistRequest) (*pb.SuccessResponse, error) {
	err := h.playlistUC.DeletePlaylist(ctx, req.PlaylistId)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "failed to delete playlist: %v", err)
	}
	return &pb.SuccessResponse{Success: true}, nil
}
