package domain

import "context"

type ArtistRepository interface {
	Create(ctx context.Context, artist *Artist) (int64, error)
	GetByID(ctx context.Context, id int64) (*Artist, error)
	Search(ctx context.Context, query string, limit int32) ([]Artist, error)
}

type AlbumRepository interface {
	Create(ctx context.Context, album *Album) (int64, error)
	GetByID(ctx context.Context, id int64) (*Album, error)
	GetByArtistID(ctx context.Context, artistID int64) ([]Album, error)
	Search(ctx context.Context, query string, limit int32) ([]Album, error)
}

type SongRepository interface {
	Create(ctx context.Context, song *Song) (int64, error)
	GetByID(ctx context.Context, id int64) (*Song, error)
	Update(ctx context.Context, song *Song) error
	Delete(ctx context.Context, id int64) error
	GetByAlbumID(ctx context.Context, albumID int64) ([]Song, error)
	GetByGenre(ctx context.Context, genre string, limit int32) ([]Song, error)
	Search(ctx context.Context, query string, limit int32) ([]Song, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expirationSeconds int) error
	Get(ctx context.Context, key string, dest interface{}) error
}
