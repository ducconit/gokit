package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisSource struct {
	Client *redis.Client
	Key    string
}

func (s *RedisSource) Load(ctx context.Context) ([]byte, error) {
	if s.Client == nil {
		return nil, fmt.Errorf("config: redis missing client")
	}
	if s.Key == "" {
		return nil, fmt.Errorf("config: redis missing key")
	}

	return s.Client.Get(ctx, s.Key).Bytes()
}
