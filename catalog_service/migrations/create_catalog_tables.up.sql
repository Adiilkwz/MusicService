CREATE TABLE artists (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT
);

CREATE TABLE albums (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id BIGINT REFERENCES artists(id) ON DELETE CASCADE,
    release_year INT
);

CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id BIGINT REFERENCES artists(id) ON DELETE CASCADE,
    album_id BIGINT REFERENCES albums(id) ON DELETE SET NULL,
    duration_seconds INT NOT NULL,
    cover_image_url VARCHAR(512)
);

CREATE INDEX idx_songs_title ON songs(title);
CREATE INDEX idx_artists_name ON artists(name);