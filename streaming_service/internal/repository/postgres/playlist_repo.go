package postgres

import (
	"context"

	"streaming_service/internal/domain"
)

type playlistRepo struct {
	db *DB
}

func NewPlaylistRepository(db *DB) domain.PlaylistRepository {
	return &playlistRepo{db: db}
}

func (r *playlistRepo) Create(ctx context.Context, userID int64, title string) (int64, error) {
	query := `INSERT INTO playlists (user_id, title) VALUES ($1, $2) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, userID, title).Scan(&id)
	return id, err
}

func (r *playlistRepo) GetByID(ctx context.Context, playlistID int64) (*domain.Playlist, error) {
	query := `SELECT id, user_id, title FROM playlists WHERE id = $1`
	var p domain.Playlist
	err := r.db.QueryRowContext(ctx, query, playlistID).Scan(&p.ID, &p.UserID, &p.Title)
	if err != nil {
		return nil, err
	}

	songQuery := `SELECT song_id FROM playlist_songs WHERE playlist_id = $1 ORDER BY position`
	rows, err := r.db.QueryContext(ctx, songQuery, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var songID int64
		if err := rows.Scan(&songID); err != nil {
			return nil, err
		}
		p.SongIDs = append(p.SongIDs, songID)
	}
	return &p, nil
}

func (r *playlistRepo) AddSong(ctx context.Context, playlistID, songID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2)`
	if err := tx.QueryRowContext(ctx, checkQuery, playlistID, songID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	var maxPos int
	posQuery := `SELECT COALESCE(MAX(position), 0) FROM playlist_songs WHERE playlist_id = $1`
	if err := tx.QueryRowContext(ctx, posQuery, playlistID).Scan(&maxPos); err != nil {
		return err
	}

	insertQuery := `INSERT INTO playlist_songs (playlist_id, song_id, position) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, insertQuery, playlistID, songID, maxPos+1)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *playlistRepo) RemoveSong(ctx context.Context, playlistID, songID int64) error {
	query := `DELETE FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`
	_, err := r.db.ExecContext(ctx, query, playlistID, songID)
	return err
}

func (r *playlistRepo) Delete(ctx context.Context, playlistID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM playlist_songs WHERE playlist_id = $1`, playlistID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM playlists WHERE id = $1`, playlistID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
