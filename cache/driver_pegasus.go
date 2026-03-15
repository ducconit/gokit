package cache

import (
	"context"
	"fmt"

	"github.com/eko/gocache/lib/v4/store"
	pegasus_store "github.com/eko/gocache/store/pegasus/v4"
)

func newPegasusStore(ctx context.Context, cfg Config) (store.StoreInterface, error) {
	if len(cfg.Pegasus.MetaServers) == 0 {
		return nil, fmt.Errorf("pegasus: missing meta servers")
	}

	st, err := pegasus_store.NewPegasus(ctx, &pegasus_store.OptionsPegasus{
		MetaServers: cfg.Pegasus.MetaServers,
	})
	if err != nil {
		return nil, err
	}
	return st, nil
}
