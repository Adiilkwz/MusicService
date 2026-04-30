package postgres

import (
	"context"
	"catalog_service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type albumRepo struct {
	db *pgxpool.Pool
}

func NewAlbumRepository(db *pgxpool.Pool) domain.AlbumRepository {
	return &albumRepo{db: db}
}

func (r *albumRepo) Create(ctx context.Context, album *domain.Album) (int64, error) {
	query := `INSERT INTO albums (title, artist_id, release_year) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRow(ctx, query, album.Title, album.ArtistID, album.ReleaseYear).Scan(&id)
	return id, err
}

func (r *albumRepo) GetByID(ctx context.Context, id int64) (*domain.Album, error) {
	query := `SELECT id, title, artist_id, release_year FROM albums WHERE id = $1`
	a := &domain.Album{}
	err := r.db.QueryRow(ctx, query, id).Scan(&a.ID, &a.Title, &a.ArtistID, &a.ReleaseYear)
	return a, err
}

func (r *albumRepo) GetByArtistID(ctx context.Context, artistID int64) ([]domain.Album, error) {
	query := `SELECT id, title, artist_id, release_year FROM albums WHERE artist_id = $1`
	rows, err := r.db.Query(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []domain.Album
	for rows.Next() {
		var a domain.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.ReleaseYear); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}

func (r *albumRepo) Search(ctx context.Context, query string, limit int32) ([]domain.Album, error) {
	sqlQuery := `SELECT id, title, artist_id, release_year FROM albums WHERE title ILIKE '%' || $1 || '%' LIMIT $2`
	rows, err := r.db.Query(ctx, sqlQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []domain.Album
	for rows.Next() {
		var a domain.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.ReleaseYear); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}