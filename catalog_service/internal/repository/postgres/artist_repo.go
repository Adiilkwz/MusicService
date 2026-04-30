package postgres

import (
	"context"

	"catalog_service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type artistRepo struct {
	db *pgxpool.Pool
}

func NewArtistRepository(db *pgxpool.Pool) domain.ArtistRepository {
	return &artistRepo{db: db}
}

func (r *artistRepo) Create(ctx context.Context, artist *domain.Artist) (int64, error) {
	query := `INSERT INTO artists (name, bio) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.db.QueryRow(ctx, query, artist.Name, artist.Bio).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *artistRepo) GetByID(ctx context.Context, id int64) (*domain.Artist, error) {
	query := `SELECT id, name, bio FROM artists WHERE id = $1`

	artist := &domain.Artist{}
	err := r.db.QueryRow(ctx, query, id).Scan(&artist.ID, &artist.Name, &artist.Bio)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (r *artistRepo) Search(ctx context.Context, query string, limit int32) ([]domain.Artist, error) {
	sqlQuery := `SELECT id, name, bio FROM artists WHERE name ILIKE '%' || $1 || '%' LIMIT $2`

	rows, err := r.db.Query(ctx, sqlQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []domain.Artist
	for rows.Next() {
		var a domain.Artist
		if err := rows.Scan(&a.ID, &a.Name, &a.Bio); err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}

	return artists, nil
}
