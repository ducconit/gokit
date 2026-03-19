package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ducconit/gokit/id"
	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	SessionID        string
}

type Manager struct {
	cfg   Config
	store SessionStore
	now   func() time.Time
}

func New(cfg Config, store SessionStore) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("auth: missing session store")
	}

	// If it's a SQLStore, inject the execution logic from config or use defaults if DB is provided
	if sqlStore, ok := store.(*SQLStore); ok {
		if cfg.SQLExec != nil {
			sqlStore.SQLExec = cfg.SQLExec
		} else if sqlStore.DB != nil {
			sqlStore.SQLExec = func(ctx context.Context, query string, args ...any) (sql.Result, error) {
				return sqlStore.DB.ExecContext(ctx, query, args...)
			}
		}

		if cfg.SQLQueryRow != nil {
			sqlStore.SQLQueryRow = cfg.SQLQueryRow
		} else if sqlStore.DB != nil {
			sqlStore.SQLQueryRow = func(ctx context.Context, query string, args ...any) *sql.Row {
				return sqlStore.DB.QueryRowContext(ctx, query, args...)
			}
		}
	}

	if len(cfg.HMACSecret) == 0 {
		return nil, fmt.Errorf("auth: missing HMAC secret")
	}
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 15 * time.Minute
	}
	if cfg.RefreshTTL <= 0 {
		cfg.RefreshTTL = 30 * 24 * time.Hour
	}

	return &Manager{
		cfg:   cfg,
		store: store,
		now:   time.Now,
	}, nil
}

func (m *Manager) Issue(ctx context.Context, subjectID, subjectType string, metadata map[string]string) (TokenPair, error) {
	now := m.now()
	sid := id.New()

	refreshExp := now.Add(m.cfg.RefreshTTL)
	refresh, err := m.signRefreshToken(now, refreshExp, subjectID, subjectType, sid)
	if err != nil {
		return TokenPair{}, err
	}

	s := Session{
		ID:              sid,
		SubjectID:       subjectID,
		SubjectType:     subjectType,
		CreatedAt:       now,
		ExpiresAt:       refreshExp,
		RefreshTokenSum: sha256.Sum256([]byte(refresh)),
		Metadata:        metadata,
	}
	if err := m.store.Put(ctx, s); err != nil {
		return TokenPair{}, err
	}

	accessExp := now.Add(m.cfg.AccessTTL)
	access, err := m.signAccessToken(now, accessExp, subjectID, subjectType, sid)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
		SessionID:        sid,
	}, nil
}

func (m *Manager) UpdateSessionMetadata(ctx context.Context, sessionID string, metadata map[string]string) error {
	return m.store.UpdateMetadata(ctx, sessionID, metadata)
}

func (m *Manager) VerifyAccess(ctx context.Context, token string) (AccessClaims, error) {
	var claims AccessClaims

	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("auth: unexpected signing method: %s", token.Method.Alg())
		}
		return m.cfg.HMACSecret, nil
	}, jwt.WithAudience(m.cfg.Audience...), jwt.WithIssuer(m.cfg.Issuer))
	if err != nil {
		return AccessClaims{}, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	if !t.Valid {
		return AccessClaims{}, ErrUnauthorized
	}

	// 1. Check if subject is banned
	banned, err := m.store.IsSubjectBanned(ctx, claims.SubjectID)
	if err != nil {
		return AccessClaims{}, err
	}
	if banned {
		return AccessClaims{}, ErrForbidden
	}

	// 2. Check if session exists and is valid
	s, err := m.store.Get(ctx, claims.SessionID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return AccessClaims{}, ErrUnauthorized
		}
		return AccessClaims{}, err
	}
	if m.now().After(s.ExpiresAt) {
		return AccessClaims{}, ErrUnauthorized
	}

	return claims, nil
}

type RefreshOptions struct {
	Rotate bool
}

func (m *Manager) Refresh(ctx context.Context, refreshToken string, opts RefreshOptions) (TokenPair, error) {
	var claims RefreshClaims

	t, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("auth: unexpected signing method: %s", token.Method.Alg())
		}
		return m.cfg.HMACSecret, nil
	}, jwt.WithAudience(m.cfg.Audience...), jwt.WithIssuer(m.cfg.Issuer))
	if err != nil || !t.Valid {
		return TokenPair{}, ErrUnauthorized
	}

	s, err := m.store.Get(ctx, claims.SessionID)
	if err != nil {
		return TokenPair{}, ErrUnauthorized
	}
	if m.now().After(s.ExpiresAt) {
		return TokenPair{}, ErrUnauthorized
	}
	if s.RefreshTokenSum != sha256.Sum256([]byte(refreshToken)) {
		return TokenPair{}, ErrUnauthorized
	}

	now := m.now()
	var newRefresh string
	newRefreshExp := s.ExpiresAt

	if opts.Rotate {
		newRefreshExp = now.Add(m.cfg.RefreshTTL)
		newRefresh, err = m.signRefreshToken(now, newRefreshExp, s.SubjectID, s.SubjectType, s.ID)
		if err != nil {
			return TokenPair{}, err
		}
		s.ExpiresAt = newRefreshExp
		s.RefreshTokenSum = sha256.Sum256([]byte(newRefresh))
		if err := m.store.Put(ctx, s); err != nil {
			return TokenPair{}, err
		}
	} else {
		newRefresh = refreshToken
	}

	newAccessExp := now.Add(m.cfg.AccessTTL)
	newAccess, err := m.signAccessToken(now, newAccessExp, s.SubjectID, s.SubjectType, s.ID)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      newAccess,
		RefreshToken:     newRefresh,
		AccessExpiresAt:  newAccessExp,
		RefreshExpiresAt: newRefreshExp,
		SessionID:        s.ID,
	}, nil
}

func (m *Manager) Logout(ctx context.Context, sessionID string) error {
	return m.store.Delete(ctx, sessionID)
}

func (m *Manager) LogoutAll(ctx context.Context, subjectID string) error {
	return m.store.DeleteBySubject(ctx, subjectID)
}

func (m *Manager) BanSubject(ctx context.Context, subjectID string, until time.Time) error {
	return m.store.BanSubject(ctx, subjectID, until)
}

func (m *Manager) signAccessToken(issuedAt, expiresAt time.Time, subjectID, subjectType, sessionID string) (string, error) {
	claims := AccessClaims{
		SubjectID:   subjectID,
		SubjectType: subjectType,
		SessionID:   sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Audience:  jwt.ClaimStrings(m.cfg.Audience),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   subjectID,
			ID:        id.New(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.cfg.HMACSecret)
}

func (m *Manager) signRefreshToken(issuedAt, expiresAt time.Time, subjectID, subjectType, sessionID string) (string, error) {
	claims := RefreshClaims{
		SubjectID:   subjectID,
		SubjectType: subjectType,
		SessionID:   sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Audience:  jwt.ClaimStrings(m.cfg.Audience),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   subjectID,
			ID:        id.New(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.cfg.HMACSecret)
}
