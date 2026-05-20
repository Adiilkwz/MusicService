package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"streaming_service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type cacheRepo struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) domain.CacheRepository {
	return &cacheRepo{client: client}
}

func (r *cacheRepo) Set(ctx context.Context, key string, value interface{}, expirationSeconds int) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, b, time.Duration(expirationSeconds)*time.Second).Err()
}

func (r *cacheRepo) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
