package postgres

import (
	"context"

	"streaming_service/internal/domain"
)

type likeRepo struct {
	db *DB
}

func NewLikeRepository(db *DB) domain.LikeRepository {
	return &likeRepo{db: db}
}

func (r *likeRepo) Like(ctx context.Context, userID, songID int64) error {
	query := `INSERT INTO likes (user_id, song_id) VALUES ($1, $2) ON CONFLICT (user_id, song_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, songID)
	return err
}

func (r *likeRepo) Unlike(ctx context.Context, userID, songID int64) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND song_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, songID)
	return err
}

func (r *likeRepo) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]int64, error) {
	query := `SELECT song_id FROM likes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songIDs []int64
	for rows.Next() {
		var songID int64
		if err := rows.Scan(&songID); err != nil {
			return nil, err
		}
		songIDs = append(songIDs, songID)
	}
	return songIDs, nil
}
