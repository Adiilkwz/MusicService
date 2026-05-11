package grpc

import (
	"catalog_service/internal/domain"

	"github.com/Adiilkwz/music-grpc-go/catalog"
)

type Server struct {
	catalog.UnimplementedCatalogServiceServer
	artistUC domain.ArtistUsecase
	albumUC  domain.AlbumUsecase
	songUC   domain.SongUsecase
	searchUC domain.SearchUsecase
}

func NewServer(a domain.ArtistUsecase, al domain.AlbumUsecase, s domain.SongUsecase, se domain.SearchUsecase) *Server {
	return &Server{
		artistUC: a,
		albumUC:  al,
		songUC:   s,
		searchUC: se,
	}
}
