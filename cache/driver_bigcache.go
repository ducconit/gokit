package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
)

func newBigCacheStore(cfg Config) (store.StoreInterface, error) {
	life := cfg.BigCache.LifeWindow
	if life <= 0 {
		life = 5 * time.Minute
	}

	bc := bigcache.DefaultConfig(life)
	if cfg.BigCache.CleanWindow > 0 {
		bc.CleanWindow = cfg.BigCache.CleanWindow
	}
	if cfg.BigCache.Shards > 0 {
		bc.Shards = cfg.BigCache.Shards
	}
	if cfg.BigCache.MaxEntriesInWindow > 0 {
		bc.MaxEntriesInWindow = cfg.BigCache.MaxEntriesInWindow
	}
	if cfg.BigCache.MaxEntrySize > 0 {
		bc.MaxEntrySize = cfg.BigCache.MaxEntrySize
	}
	if cfg.BigCache.HardMaxCacheSize > 0 {
		bc.HardMaxCacheSize = cfg.BigCache.HardMaxCacheSize
	}
	bc.Verbose = cfg.BigCache.Verbose

	client, err := bigcache.NewBigCache(bc)
	if err != nil {
		return nil, err
	}

	return bigcache_store.NewBigcache(client), nil
}
