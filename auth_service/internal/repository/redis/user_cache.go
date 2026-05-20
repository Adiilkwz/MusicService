package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"auth_service/internal/domain"

	redis_client "github.com/redis/go-redis/v9"
)

type userCache struct {
	client *redis_client.Client
	ttl    time.Duration
}

func NewUserCache(client *redis_client.Client, ttl time.Duration) domain.UserCache {
	return &userCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *userCache) SetUser(ctx context.Context, user *domain.User) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d", user.ID)
	return c.client.Set(ctx, key, bytes, c.ttl).Err()
}

func (c *userCache) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *userCache) DeleteUser(ctx context.Context, id int64) error {
	key := fmt.Sprintf("user:%d", id)
	return c.client.Del(ctx, key).Err()
}
