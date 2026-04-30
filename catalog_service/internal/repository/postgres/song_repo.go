package postgres

import (
	"context"
	"catalog_service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type songRepo struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) domain.SongRepository {
	return &songRepo{db: db}
}

func (r *songRepo) Create(ctx context.Context, s *domain.Song) (int64, error) {
	query := `INSERT INTO songs (title, artist_id, album_id, duration_seconds, genre) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := r.db.QueryRow(ctx, query, s.Title, s.ArtistID, s.AlbumID, s.DurationSeconds, s.Genre).Scan(&id)
	return id, err
}

func (r *songRepo) GetByID(ctx context.Context, id int64) (*domain.Song, error) {
	query := `SELECT id, title, artist_id, album_id, duration_seconds, genre FROM songs WHERE id = $1`
	s := &domain.Song{}
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.Title, &s.ArtistID, &s.AlbumID, &s.DurationSeconds, &s.Genre)
	return s, err
}

func (r *songRepo) Update(ctx context.Context, s *domain.Song) error {
	query := `UPDATE songs SET title = $1, genre = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, s.Title, s.Genre, s.ID)
	return err
}

func (r *songRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *songRepo) GetByGenre(ctx context.Context, genre string, limit int32) ([]domain.Song, error) {
	query := `SELECT id, title, artist_id, album_id, duration_seconds, genre FROM songs WHERE genre = $1 LIMIT $2`
	rows, err := r.db.Query(ctx, query, genre, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []domain.Song
	for rows.Next() {
		var s domain.Song
		if err := rows.Scan(&s.ID, &s.Title, &s.ArtistID, &s.AlbumID, &s.DurationSeconds, &s.Genre); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

func (r *songRepo) Search(ctx context.Context, query string, limit int32) ([]domain.Song, error) {
	sqlQuery := `SELECT id, title, artist_id, album_id, duration_seconds, genre FROM songs WHERE title ILIKE '%' || $1 || '%' LIMIT $2`
	rows, err := r.db.Query(ctx, sqlQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []domain.Song
	for rows.Next() {
		var s domain.Song
		if err := rows.Scan(&s.ID, &s.Title, &s.ArtistID, &s.AlbumID, &s.DurationSeconds, &s.Genre); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}