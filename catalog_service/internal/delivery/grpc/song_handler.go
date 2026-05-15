package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ensureHasRole(ctx context.Context, allowedRoles ...string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not found")
	}

	roles := md.Get("user_role")
	if len(roles) == 0 {
		return status.Error(codes.Unauthenticated, "user_role missing in metadata")
	}

	userRole := roles[0]
	for _, allowed := range allowedRoles {
		if userRole == allowed {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, "role '%s' does not have permission for this action", userRole)
}

func (s *Server) CreateSong(ctx context.Context, req *catalog.CreateSongRequest) (*catalog.CreateSongResponse, error) {
	if err := ensureHasRole(ctx, "admin", "artist"); err != nil {
		return nil, err
	}

	id, err := s.songUC.CreateSong(ctx, req.GetAlbumId(), req.GetTitle(), req.GetDuration(), req.GetGenre())
	if err != nil {
		return nil, err
	}
	return &catalog.CreateSongResponse{Id: id}, nil
}

func (s *Server) GetSong(ctx context.Context, req *catalog.GetSongRequest) (*catalog.GetSongResponse, error) {
	song, err := s.songUC.GetSong(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	album, _ := s.albumUC.GetAlbum(ctx, song.AlbumID)
	artist, _ := s.artistUC.GetArtist(ctx, album.ArtistID)

	return &catalog.GetSongResponse{
		Song: &catalog.Song{
			Id:       song.ID,
			AlbumId:  song.AlbumID,
			Title:    song.Title,
			Duration: song.DurationSeconds,
			Genre:    song.Genre,
		},
		ArtistName: artist.Name,
		AlbumTitle: album.Title,
	}, nil
}

func (s *Server) UpdateSong(ctx context.Context, req *catalog.UpdateSongRequest) (*catalog.SuccessResponse, error) {
	if err := ensureHasRole(ctx, "admin", "artist"); err != nil {
		return &catalog.SuccessResponse{Success: false}, err
	}

	err := s.songUC.UpdateSong(ctx, req.GetId(), req.GetTitle(), req.GetGenre())
	if err != nil {
		return &catalog.SuccessResponse{Success: false}, err
	}
	return &catalog.SuccessResponse{Success: true}, nil
}

func (s *Server) DeleteSong(ctx context.Context, req *catalog.DeleteSongRequest) (*catalog.SuccessResponse, error) {
	if err := ensureHasRole(ctx, "admin"); err != nil {
		return &catalog.SuccessResponse{Success: false}, err
	}

	err := s.songUC.DeleteSong(ctx, req.GetId())
	if err != nil {
		return &catalog.SuccessResponse{Success: false}, err
	}
	return &catalog.SuccessResponse{Success: true}, nil
}

func (s *Server) GetSongsByGenre(ctx context.Context, req *catalog.GetSongsByGenreRequest) (*catalog.GetSongsByGenreResponse, error) {
	songs, err := s.songUC.GetSongsByGenre(ctx, req.GetGenre(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	var pbSongs []*catalog.Song
	for _, so := range songs {
		pbSongs = append(pbSongs, &catalog.Song{
			Id:       so.ID,
			AlbumId:  so.AlbumID,
			Title:    so.Title,
			Duration: so.DurationSeconds,
			Genre:    so.Genre,
		})
	}
	return &catalog.GetSongsByGenreResponse{Songs: pbSongs}, nil
}
