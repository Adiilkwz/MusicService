package postgres

import (
	"context"

	"streaming_service/internal/domain"
)

type trendingRepo struct {
	db *DB
}

func NewTrendingRepository(db *DB) domain.TrendingRepository {
	return &trendingRepo{db: db}
}

func (r *trendingRepo) GetTrending(ctx context.Context, limit int) ([]domain.TrendingItem, error) {
	query := `SELECT song_id, play_count FROM trending_songs ORDER BY play_count DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.TrendingItem
	for rows.Next() {
		var item domain.TrendingItem
		if err := rows.Scan(&item.SongID, &item.PlayCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *trendingRepo) IncrementPlayCount(ctx context.Context, songID int64) error {
	query := `INSERT INTO trending_songs (song_id, play_count) VALUES ($1, 1)
			  ON CONFLICT (song_id) DO UPDATE SET play_count = trending_songs.play_count + 1`
	_, err := r.db.ExecContext(ctx, query, songID)
	return err
}
