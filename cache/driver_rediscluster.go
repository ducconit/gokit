package cache

import (
	"fmt"

	"github.com/eko/gocache/lib/v4/store"
	rediscluster_store "github.com/eko/gocache/store/rediscluster/v4"
	"github.com/redis/go-redis/v9"
)

func newRedisClusterStore(cfg Config) (store.StoreInterface, error) {
	if len(cfg.RedisCluster.Addrs) == 0 {
		return nil, fmt.Errorf("rediscluster: missing addrs")
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.RedisCluster.Addrs,
		Username: cfg.RedisCluster.Username,
		Password: cfg.RedisCluster.Password,
	})

	if cfg.DefaultExpiration > 0 {
		return rediscluster_store.NewRedisCluster(client, store.WithExpiration(cfg.DefaultExpiration)), nil
	}
	return rediscluster_store.NewRedisCluster(client), nil
}
