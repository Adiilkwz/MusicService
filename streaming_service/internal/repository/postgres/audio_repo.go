package postgres

import (
	"fmt"
	"strconv"

	"streaming_service/internal/domain"
)

type audioRepo struct {
	audioDir string
}

func NewAudioRepository(audioDir string) domain.AudioRepository {
	return &audioRepo{audioDir: audioDir}
}

func (r *audioRepo) GetAudioPath(songID int64) string {
	return fmt.Sprintf("%s/song_%s.mp3", r.audioDir, strconv.FormatInt(songID, 10))
}

func (r *audioRepo) ReadChunk(songID int64, offset, size int64) ([]byte, error) {
	return make([]byte, size), nil
}
