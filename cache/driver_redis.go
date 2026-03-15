package cache

import (
	"fmt"

	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

func newRedisStore(cfg Config) (store.StoreInterface, error) {
	if cfg.Redis.Addr == "" {
		return nil, fmt.Errorf("redis: missing addr")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if cfg.DefaultExpiration > 0 {
		return redis_store.NewRedis(client, store.WithExpiration(cfg.DefaultExpiration)), nil
	}
	return redis_store.NewRedis(client), nil
}
