\c catalog_db

CREATE TABLE IF NOT EXISTS artists (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT
);

CREATE TABLE IF NOT EXISTS albums (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id BIGINT REFERENCES artists(id) ON DELETE CASCADE,
    release_year INT
);

CREATE TABLE IF NOT EXISTS songs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id BIGINT REFERENCES artists(id) ON DELETE CASCADE,
    album_id BIGINT REFERENCES albums(id) ON DELETE SET NULL,
    duration_seconds INT NOT NULL,
    genre VARCHAR(100),
    cover_image_url VARCHAR(512)
);

CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_artists_name ON artists(name);
CREATE INDEX IF NOT EXISTS idx_songs_genre ON songs(genre);
