package cache

import (
	"context"
	"fmt"
	"time"

	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
)

var (
	ErrUnknownDriver = fmt.Errorf("unknown cache driver")
)

type Manager[T any] struct {
	c                 gocache.CacheInterface[T]
	defaultExpiration time.Duration
}

func New[T any](ctx context.Context, cfg Config) (*Manager[T], error) {
	st, err := newStore[T](ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Manager[T]{
		c:                 gocache.New[T](st),
		defaultExpiration: cfg.DefaultExpiration,
	}, nil
}

func (m *Manager[T]) Get(ctx context.Context, key any) (T, error) {
	return m.c.Get(ctx, key)
}

func (m *Manager[T]) Set(ctx context.Context, key any, value T, opts ...store.Option) error {
	if m.defaultExpiration > 0 {
		opts = append([]store.Option{store.WithExpiration(m.defaultExpiration)}, opts...)
	}
	return m.c.Set(ctx, key, value, opts...)
}

func (m *Manager[T]) Delete(ctx context.Context, key any) error {
	return m.c.Delete(ctx, key)
}

func (m *Manager[T]) Clear(ctx context.Context) error {
	return m.c.Clear(ctx)
}

func (m *Manager[T]) Invalidate(ctx context.Context, opts ...store.InvalidateOption) error {
	return m.c.Invalidate(ctx, opts...)
}

func newStore[T any](ctx context.Context, cfg Config) (store.StoreInterface, error) {
	switch cfg.Driver {
	case DriverBigCache:
		return newBigCacheStore(cfg)
	case DriverFreeCache:
		return newFreeCacheStore(cfg)
	case DriverGoCache:
		return newGoCacheStore(cfg)
	case DriverRistretto:
		return newRistrettoStore(cfg)
	case DriverMemcache:
		return newMemcacheStore(cfg)
	case DriverRedis:
		return newRedisStore(cfg)
	case DriverRedisCluster:
		return newRedisClusterStore(cfg)
	case DriverRueidis:
		return newRueidisStore(cfg)
	case DriverHazelcast:
		return newHazelcastStore(ctx, cfg)
	case DriverPegasus:
		return newPegasusStore(ctx, cfg)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownDriver, cfg.Driver)
	}
}
