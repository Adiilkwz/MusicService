package grpc

import (
	"context"

	"github.com/Adiilkwz/music-grpc-go/catalog"
)

func (s *Server) SearchCatalog(ctx context.Context, req *catalog.SearchCatalogRequest) (*catalog.SearchCatalogResponse, error) {
	res, err := s.searchUC.SearchCatalog(ctx, req.GetQuery(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	resp := &catalog.SearchCatalogResponse{}

	for _, a := range res.Artists {
		resp.Artists = append(resp.Artists, &catalog.Artist{Id: a.ID, Name: a.Name, Bio: a.Bio})
	}
	for _, al := range res.Albums {
		resp.Albums = append(resp.Albums, &catalog.Album{Id: al.ID, ArtistId: al.ArtistID, Title: al.Title, ReleaseYear: al.ReleaseYear})
	}
	for _, so := range res.Songs {
		resp.Songs = append(resp.Songs, &catalog.Song{Id: so.ID, AlbumId: so.AlbumID, Title: so.Title, Duration: so.DurationSeconds, Genre: so.Genre})
	}

	return resp, nil
}
