package cache

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/eko/gocache/lib/v4/store"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

func newRistrettoStore(cfg Config) (store.StoreInterface, error) {
	numCounters := cfg.Ristretto.NumCounters
	if numCounters <= 0 {
		numCounters = 10_000
	}
	maxCost := cfg.Ristretto.MaxCost
	if maxCost <= 0 {
		maxCost = 1_000_000
	}
	bufferItems := cfg.Ristretto.BufferItems
	if bufferItems <= 0 {
		bufferItems = 64
	}

	client, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		return nil, err
	}

	return ristretto_store.NewRistretto[string, any](client), nil
}
