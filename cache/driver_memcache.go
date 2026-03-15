package cache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/eko/gocache/lib/v4/store"
	memcache_store "github.com/eko/gocache/store/memcache/v4"
)

func newMemcacheStore(cfg Config) (store.StoreInterface, error) {
	if len(cfg.Memcache.Servers) == 0 {
		return nil, fmt.Errorf("memcache: missing servers")
	}

	client := memcache.New(cfg.Memcache.Servers...)
	if cfg.DefaultExpiration > 0 {
		return memcache_store.NewMemcache(client, store.WithExpiration(cfg.DefaultExpiration)), nil
	}
	return memcache_store.NewMemcache(client), nil
}
