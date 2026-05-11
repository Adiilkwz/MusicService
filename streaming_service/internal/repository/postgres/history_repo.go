package postgres

import (
	"context"
	"time"

	"streaming_service/internal/domain"
)

type historyRepo struct {
	db *DB
}

func NewHistoryRepository(db *DB) domain.HistoryRepository {
	return &historyRepo{db: db}
}

func (r *historyRepo) Create(ctx context.Context, userID, songID int64) error {
	query := `INSERT INTO play_history (user_id, song_id, played_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, userID, songID, time.Now())
	return err
}

func (r *historyRepo) GetByUserID(ctx context.Context, userID int64, limit int) ([]domain.History, error) {
	query := `SELECT id, user_id, song_id, played_at FROM play_history 
			  WHERE user_id = $1 ORDER BY played_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.History
	for rows.Next() {
		var h domain.History
		if err := rows.Scan(&h.ID, &h.UserID, &h.SongID, &h.PlayedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
