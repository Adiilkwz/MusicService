\c streaming_db

CREATE TABLE IF NOT EXISTS play_history (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    song_id BIGINT NOT NULL,
    played_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS playlists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS playlist_songs (
    playlist_id BIGINT REFERENCES playlists(id) ON DELETE CASCADE,
    song_id BIGINT NOT NULL,
    position INTEGER DEFAULT 0,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (playlist_id, song_id)
);

CREATE TABLE IF NOT EXISTS likes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    song_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, song_id)
);

CREATE TABLE IF NOT EXISTS trending_songs (
    song_id BIGINT PRIMARY KEY,
    play_count INTEGER DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_play_history_user_id ON play_history(user_id);
CREATE INDEX IF NOT EXISTS idx_play_history_song_id ON play_history(song_id);
CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
CREATE INDEX IF NOT EXISTS idx_trending_play_count ON trending_songs(play_count DESC);
