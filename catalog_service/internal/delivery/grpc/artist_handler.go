package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/catalog"
)

func (s *Server) CreateArtist(ctx context.Context, req *catalog.CreateArtistRequest) (*catalog.CreateArtistResponse, error) {
	if err := ensureHasRole(ctx, "admin", "artist"); err != nil {
		return nil, err
	}

	id, err := s.artistUC.CreateArtist(ctx, req.GetName(), req.GetBio())
	if err != nil {
		return nil, err
	}
	return &catalog.CreateArtistResponse{Id: id}, nil
}

func (s *Server) GetArtist(ctx context.Context, req *catalog.GetArtistRequest) (*catalog.GetArtistResponse, error) {
	a, err := s.artistUC.GetArtist(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &catalog.GetArtistResponse{
		Artist: &catalog.Artist{Id: a.ID, Name: a.Name, Bio: a.Bio},
	}, nil
}

func (s *Server) GetAlbumsByArtist(ctx context.Context, req *catalog.GetAlbumsByArtistRequest) (*catalog.GetAlbumsByArtistResponse, error) {
	albums, err := s.artistUC.GetAlbumsByArtist(ctx, req.GetArtistId())
	if err != nil {
		return nil, err
	}

	var pbAlbums []*catalog.Album
	for _, al := range albums {
		pbAlbums = append(pbAlbums, &catalog.Album{
			Id:          al.ID,
			ArtistId:    al.ArtistID,
			Title:       al.Title,
			ReleaseYear: al.ReleaseYear,
		})
	}
	return &catalog.GetAlbumsByArtistResponse{Albums: pbAlbums}, nil
}
