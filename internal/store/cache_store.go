package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
	SetResetCode(email, code string, expiration time.Duration) error
	GetResetCode(email string) (string, error)
	DeleteResetCode(email string) error
	Close() error
}

type RedisCacheStore struct {
	client *redis.Client
}

func NewRedisCacheStore(client *redis.Client) *RedisCacheStore {
	return &RedisCacheStore{
		client: client,
	}
}

func (r *RedisCacheStore) SetResetCode(email, code string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("reset_code:%s", email)
	return r.client.Set(ctx, key, code, expiration).Err()
}

func (r *RedisCacheStore) GetResetCode(email string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("reset_code:%s", email)
	code, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("reset code not found or expired")
	}
	return code, err
}

func (r *RedisCacheStore) DeleteResetCode(email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("reset_code:%s", email)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCacheStore) Close() error {
	return r.client.Close()
}
