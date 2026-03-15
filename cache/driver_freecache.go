package cache

import (
	"github.com/coocood/freecache"
	"github.com/eko/gocache/lib/v4/store"
	freecache_store "github.com/eko/gocache/store/freecache/v4"
)

func newFreeCacheStore(cfg Config) (store.StoreInterface, error) {
	size := cfg.FreeCache.Size
	if size <= 0 {
		size = 100 * 1024 * 1024
	}
	client := freecache.NewCache(size)
	return freecache_store.NewFreecache(client), nil
}
