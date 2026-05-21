-- Run in pgAdmin against the EXISTING catalog_db (do not create a new database).
-- Connection: host localhost, port 5435, user postgres, database catalog_db
--
-- Audio files directory (host): streaming_service/audio_files/
-- Docker mount: ./streaming_service/audio_files -> /app/audio_files
--
-- Add MP3s to that folder. Streaming opens: /app/audio_files/song_<id>.mp3
-- After INSERT, rename/copy files to match song ids, for example:
--   track1.mp3 -> song_1.mp3
--   track2.mp3 -> song_2.mp3
--   track3.mp3 -> song_3.mp3

\c catalog_db

ALTER TABLE songs ADD COLUMN IF NOT EXISTS audio_url VARCHAR(512);

BEGIN;

INSERT INTO artists (name, bio)
SELECT 'Demo Artist', 'Sample artist for local streaming tests'
WHERE NOT EXISTS (SELECT 1 FROM artists WHERE name = 'Demo Artist');

INSERT INTO albums (title, artist_id, release_year)
SELECT 'Demo Album', a.id, 2024
FROM artists a
WHERE a.name = 'Demo Artist'
  AND NOT EXISTS (
    SELECT 1 FROM albums al WHERE al.title = 'Demo Album' AND al.artist_id = a.id
  );

INSERT INTO songs (title, artist_id, album_id, duration_seconds, genre, cover_image_url, audio_url)
SELECT v.title, a.id, al.id, v.duration_seconds, v.genre, v.cover_image_url, v.audio_url
FROM artists a
JOIN albums al ON al.artist_id = a.id AND al.title = 'Demo Album'
CROSS JOIN (
  VALUES
    ('Track One',   210, 'Pop',  NULL::varchar, 'audio_files/track1.mp3'),
    ('Track Two',   185, 'Rock', NULL::varchar, 'audio_files/track2.mp3'),
    ('Track Three', 240, 'Jazz', NULL::varchar, 'audio_files/track3.mp3')
) AS v(title, duration_seconds, genre, cover_image_url, audio_url)
WHERE a.name = 'Demo Artist'
  AND NOT EXISTS (
    SELECT 1 FROM songs s WHERE s.album_id = al.id AND s.title = v.title
  );

COMMIT;

-- Verify:
-- SELECT s.id, s.title, s.audio_url FROM songs s ORDER BY s.id;
