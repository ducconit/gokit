package cache

import (
	"fmt"

	"github.com/eko/gocache/lib/v4/store"
	rueidis_store "github.com/eko/gocache/store/rueidis/v4"
	"github.com/redis/rueidis"
)

func newRueidisStore(cfg Config) (store.StoreInterface, error) {
	if len(cfg.Rueidis.InitAddrs) == 0 {
		return nil, fmt.Errorf("rueidis: missing init addrs")
	}

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: cfg.Rueidis.InitAddrs,
	})
	if err != nil {
		return nil, err
	}

	if cfg.DefaultExpiration > 0 {
		return rueidis_store.NewRueidis(client, store.WithExpiration(cfg.DefaultExpiration)), nil
	}
	return rueidis_store.NewRueidis(client), nil
}
