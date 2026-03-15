package cache

import (
	"context"
	"testing"
	"time"
)

func TestBigCacheSetGet(t *testing.T) {
	ctx := context.Background()

	m, err := New[[]byte](ctx, Config{
		Driver:            DriverBigCache,
		DefaultExpiration: 1 * time.Minute,
		BigCache: BigCacheConfig{
			LifeWindow: 5 * time.Minute,
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := m.Set(ctx, "k1", []byte("v1")); err != nil {
		t.Fatalf("Set: %v", err)
	}

	v, err := m.Get(ctx, "k1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(v) != "v1" {
		t.Fatalf("expected v1, got %s", string(v))
	}
}
