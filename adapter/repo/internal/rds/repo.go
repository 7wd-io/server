package rds

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type R struct {
	Client *redis.Client
}

func (dst R) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	v, err := json.Marshal(value)

	if err != nil {
		return nil
	}

	return dst.Client.Set(ctx, key, v, ttl).Err()
}

func (dst R) Get(ctx context.Context, key string, dest interface{}) error {
	v, err := dst.Client.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(v, dest)
}