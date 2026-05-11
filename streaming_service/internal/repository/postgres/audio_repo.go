package postgres

import (
	"streaming_service/internal/domain"
)

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
	return make([]byte, size), nil
}
