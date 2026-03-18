package features

import (
	"context"
	"fmt"
	"sync"
)

// MemoryStore is an in-memory storage for feature flags.
type MemoryStore struct {
	features map[string]*Feature
	mu       sync.RWMutex
}

// NewMemoryStore creates a new in-memory feature flag store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		features: make(map[string]*Feature),
	}
}

// Get retrieves a feature flag from the store.
func (s *MemoryStore) Get(ctx context.Context, key string) (*Feature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.features[key]
	if !ok {
		return nil, fmt.Errorf("feature not found: %s", key)
	}
	return f, nil
}

// All returns all feature flags in the store.
func (s *MemoryStore) All(ctx context.Context) (map[string]*Feature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string]*Feature, len(s.features))
	for k, v := range s.features {
		res[k] = v
	}
	return res, nil
}

// Set adds or updates a feature flag in the store.
func (s *MemoryStore) Set(ctx context.Context, feature *Feature) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.features[feature.Key] = feature
	return nil
}
