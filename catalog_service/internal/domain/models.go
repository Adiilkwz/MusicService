package domain

type Artist struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Bio  string `db:"bio"`
}

type Album struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	ArtistID    int64  `db:"artist_id"`
	ReleaseYear int32  `db:"release_year"`
}

type Song struct {
	ID              int64  `db:"id"`
	Title           string `db:"title"`
	ArtistID        int64  `db:"artist_id"`
	AlbumID         int64  `db:"album_id"`
	DurationSeconds int32  `db:"duration_seconds"`
	Genre           string `db:"genre"`
	CoverImageURL   string `db:"cover_image_url"`
}

type SearchResult struct {
	Artist []Artist
	Albums []Album
	Songs  []Song
}
