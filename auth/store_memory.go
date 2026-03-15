package auth

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
	userIdx  map[string]map[string]struct{}
	bans     map[string]time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: map[string]Session{},
		userIdx:  map[string]map[string]struct{}{},
		bans:     map[string]time.Time{},
	}
}

func (s *MemoryStore) Put(ctx context.Context, sess Session) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sess.ID] = sess
	if s.userIdx[sess.UserID] == nil {
		s.userIdx[sess.UserID] = map[string]struct{}{}
	}
	s.userIdx[sess.UserID][sess.ID] = struct{}{}
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, id string) (Session, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok {
		return Session{}, ErrSessionNotFound
	}
	if time.Now().After(sess.ExpiresAt) {
		return Session{}, ErrSessionNotFound
	}
	return sess, nil
}

func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[id]
	if ok {
		delete(s.sessions, id)
		if idx := s.userIdx[sess.UserID]; idx != nil {
			delete(idx, id)
			if len(idx) == 0 {
				delete(s.userIdx, sess.UserID)
			}
		}
	}
	return nil
}

func (s *MemoryStore) DeleteByUser(ctx context.Context, userID string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.userIdx[userID]
	for sid := range idx {
		delete(s.sessions, sid)
	}
	delete(s.userIdx, userID)
	return nil
}

func (s *MemoryStore) UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}
	if sess.Metadata == nil {
		sess.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		sess.Metadata[k] = v
	}
	s.sessions[id] = sess
	return nil
}

func (s *MemoryStore) BanUser(ctx context.Context, userID string, until time.Time) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bans[userID] = until
	return nil
}

func (s *MemoryStore) IsUserBanned(ctx context.Context, userID string) (bool, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	until, ok := s.bans[userID]
	if !ok {
		return false, nil
	}
	return time.Now().Before(until), nil
}
