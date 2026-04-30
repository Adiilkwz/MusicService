package domain

import "context"

type ArtistUsecase interface {
	CreateArtist(ctx context.Context, name, bio string) (int64, error)
	GetArtist(ctx context.Context, id int64) (*Artist, error)
	GetAlbumsByArtist(ctx context.Context, artistID int64) ([]Album, error)
}

type AlbumUsecase interface {
	CreateAlbum(ctx context.Context, artistID int64, title string, releaseYear int32) (int64, error)
	GetAlbum(ctx context.Context, id int64) (*Album, error)
	GetSongsByAlbum(ctx context.Context, albumID int64) ([]Song, error)
}

type SongUsecase interface {
	CreateSong(ctx context.Context, albumID int64, title string, duration int32, genre string) (int64, error)
	GetSong(ctx context.Context, id int64) (*Song, error)
	UpdateSong(ctx context.Context, id int64, title, genre string) error
	DeleteSong(ctx context.Context, id int64) error
	GetSongsByGenre(ctx context.Context, genre string, limit int32) ([]Song, error)
}

type SearchUsecase interface {
	SearchCatalog(ctx context.Context, query string, limit int32) (*SearchResult, error)
}
