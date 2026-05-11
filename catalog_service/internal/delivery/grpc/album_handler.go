package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/catalog"
)

func (s *Server) CreateAlbum(ctx context.Context, req *catalog.CreateAlbumRequest) (*catalog.CreateAlbumResponse, error) {
	id, err := s.albumUC.CreateAlbum(ctx, req.GetArtistId(), req.GetTitle(), req.GetReleaseYear())
	if err != nil {
		return nil, err
	}
	return &catalog.CreateAlbumResponse{Id: id}, nil
}

func (s *Server) GetAlbum(ctx context.Context, req *catalog.GetAlbumRequest) (*catalog.GetAlbumResponse, error) {
	album, err := s.albumUC.GetAlbum(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	artist, _ := s.artistUC.GetArtist(ctx, album.ArtistID)

	return &catalog.GetAlbumResponse{
		Album: &catalog.Album{
			Id:          album.ID,
			ArtistId:    album.ArtistID,
			Title:       album.Title,
			ReleaseYear: album.ReleaseYear,
		},
		ArtistName: artist.Name,
	}, nil
}

func (s *Server) GetSongsByAlbum(ctx context.Context, req *catalog.GetSongsByAlbumRequest) (*catalog.GetSongsByAlbumResponse, error) {
	songs, err := s.albumUC.GetSongsByAlbum(ctx, req.GetAlbumId())
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
	return &catalog.GetSongsByAlbumResponse{Songs: pbSongs}, nil
}
