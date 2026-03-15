package cache

import (
	"context"
	"fmt"

	"github.com/eko/gocache/lib/v4/store"
	hazelcast_store "github.com/eko/gocache/store/hazelcast/v4"
	"github.com/hazelcast/hazelcast-go-client"
)

func newHazelcastStore(ctx context.Context, cfg Config) (store.StoreInterface, error) {
	mapName := cfg.Hazelcast.MapName
	if mapName == "" {
		mapName = "gocache"
	}

	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		return nil, err
	}

	m, err := client.GetMap(ctx, mapName)
	if err != nil {
		return nil, fmt.Errorf("hazelcast: get map: %w", err)
	}

	return hazelcast_store.NewHazelcast(m), nil
}
