package cache

import (
	"time"

	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocacheclient "github.com/patrickmn/go-cache"
)

func newGoCacheStore(cfg Config) (store.StoreInterface, error) {
	defExp := cfg.GoCache.DefaultExpiration
	if defExp <= 0 {
		defExp = 5 * time.Minute
	}
	cleanup := cfg.GoCache.CleanupInterval
	if cleanup <= 0 {
		cleanup = 10 * time.Minute
	}
	client := gocacheclient.New(defExp, cleanup)
	return gocache_store.NewGoCache(client), nil
}
