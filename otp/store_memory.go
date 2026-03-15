package otp

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string]Record{},
	}
}

func (s *MemoryStore) Put(ctx context.Context, r Record) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key(r.Purpose, r.Recipient)] = r
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, purpose Purpose, recipient string) (Record, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.data[key(purpose, recipient)]
	if !ok {
		return Record{}, ErrNotFound
	}
	if time.Now().After(r.ExpiresAt) {
		return Record{}, ErrExpired
	}
	return r, nil
}

func (s *MemoryStore) Delete(ctx context.Context, purpose Purpose, recipient string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key(purpose, recipient))
	return nil
}

func key(purpose Purpose, recipient string) string {
	return string(purpose) + ":" + recipient
}
