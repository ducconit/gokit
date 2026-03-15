package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/ducconit/gokit/cache"
	"github.com/eko/gocache/lib/v4/store"
)

type CachingStore struct {
	Cache *cache.Manager[Record]
}

func NewCachingStore(c *cache.Manager[Record]) *CachingStore {
	return &CachingStore{
		Cache: c,
	}
}

func (s *CachingStore) Put(ctx context.Context, r Record) error {
	if s.Cache == nil {
		return fmt.Errorf("otp: cache missing manager")
	}
	ttl := time.Until(r.ExpiresAt)
	if ttl <= 0 {
		return nil // Expired
	}
	return s.Cache.Set(ctx, s.key(r.Purpose, r.Recipient), r, store.WithExpiration(ttl))
}

func (s *CachingStore) Get(ctx context.Context, purpose Purpose, recipient string) (Record, error) {
	if s.Cache == nil {
		return Record{}, fmt.Errorf("otp: cache missing manager")
	}
	r, err := s.Cache.Get(ctx, s.key(purpose, recipient))
	if err != nil {
		return Record{}, ErrNotFound
	}
	if time.Now().After(r.ExpiresAt) {
		return Record{}, ErrExpired
	}
	return r, nil
}

func (s *CachingStore) Delete(ctx context.Context, purpose Purpose, recipient string) error {
	if s.Cache == nil {
		return fmt.Errorf("otp: cache missing manager")
	}
	return s.Cache.Delete(ctx, s.key(purpose, recipient))
}

func (s *CachingStore) key(purpose Purpose, recipient string) string {
	return "otp:" + string(purpose) + ":" + recipient
}
