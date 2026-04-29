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
	// Use transaction to ensure consistency
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if song already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2)`
	if err := tx.QueryRowContext(ctx, checkQuery, playlistID, songID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil // Song already in playlist
	}

	// Get max position
	var maxPos int
	posQuery := `SELECT COALESCE(MAX(position), 0) FROM playlist_songs WHERE playlist_id = $1`
	if err := tx.QueryRowContext(ctx, posQuery, playlistID).Scan(&maxPos); err != nil {
		return err
	}

	// Insert song
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

	// Delete playlist songs first
	_, err = tx.ExecContext(ctx, `DELETE FROM playlist_songs WHERE playlist_id = $1`, playlistID)
	if err != nil {
		return err
	}

	// Delete playlist
	_, err = tx.ExecContext(ctx, `DELETE FROM playlists WHERE id = $1`, playlistID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

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

type audioRepo struct {
	audioDir string
}

func NewAudioRepository(audioDir string) domain.AudioRepository {
	return &audioRepo{audioDir: audioDir}
}

func (r *audioRepo) GetAudioPath(songID int64) string {
	return r.audioDir + "/song_" + string(rune(songID)) + ".mp3"
}

func (r *audioRepo) ReadChunk(songID int64, offset, size int64) ([]byte, error) {
	// Implementation would read actual audio file
	// This is a placeholder
	return make([]byte, size), nil
}