package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ducconit/gokit/cache"
	"github.com/eko/gocache/lib/v4/store"
)

type CachingStore struct {
	Cache *cache.Manager[Session]
	Bans  *cache.Manager[time.Time]
}

func NewCachingStore(c *cache.Manager[Session], b *cache.Manager[time.Time]) *CachingStore {
	return &CachingStore{
		Cache: c,
		Bans:  b,
	}
}

func (s *CachingStore) Put(ctx context.Context, sess Session) error {
	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		return nil
	}
	return s.Cache.Set(ctx, s.sessionKey(sess.ID), sess, store.WithExpiration(ttl))
}

func (s *CachingStore) Get(ctx context.Context, id string) (Session, error) {
	sess, err := s.Cache.Get(ctx, s.sessionKey(id))
	if err != nil {
		return Session{}, ErrSessionNotFound
	}
	if time.Now().After(sess.ExpiresAt) {
		return Session{}, ErrSessionNotFound
	}
	return sess, nil
}

func (s *CachingStore) Delete(ctx context.Context, id string) error {
	return s.Cache.Delete(ctx, s.sessionKey(id))
}

func (s *CachingStore) DeleteBySubject(ctx context.Context, subjectID string) error {
	// Note: CachingStore based on key-value doesn't easily support DeleteBySubject
	// unless we maintain a secondary index (subject -> [session_ids]).
	// For now, this is a limitation of simple KV caching stores.
	return fmt.Errorf("auth: DeleteBySubject not implemented for CachingStore")
}

func (s *CachingStore) UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error {
	sess, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if sess.Metadata == nil {
		sess.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		sess.Metadata[k] = v
	}
	return s.Put(ctx, sess)
}

func (s *CachingStore) BanSubject(ctx context.Context, subjectID string, until time.Time) error {
	if s.Bans == nil {
		return fmt.Errorf("auth: bans cache missing")
	}
	ttl := time.Until(until)
	if ttl <= 0 {
		return nil
	}
	return s.Bans.Set(ctx, s.banKey(subjectID), until, store.WithExpiration(ttl))
}

func (s *CachingStore) IsSubjectBanned(ctx context.Context, subjectID string) (bool, error) {
	if s.Bans == nil {
		return false, nil
	}
	until, err := s.Bans.Get(ctx, s.banKey(subjectID))
	if err != nil {
		return false, nil
	}
	return time.Now().Before(until), nil
}

func (s *CachingStore) sessionKey(id string) string {
	return "auth:sess:" + id
}

func (s *CachingStore) banKey(subjectID string) string {
	return "auth:ban:" + subjectID
}
