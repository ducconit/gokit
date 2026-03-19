package auth

import (
	"context"
	"database/sql"
	"time"
)

type Config struct {
	Issuer     string
	Audience   []string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	HMACSecret []byte

	// SQL related configs
	SQLExec     func(ctx context.Context, query string, args ...any) (sql.Result, error)
	SQLQueryRow func(ctx context.Context, query string, args ...any) *sql.Row

	Store SessionStore
}

type Option func(*Config)

func WithIssuer(issuer string) Option {
	return func(c *Config) {
		c.Issuer = issuer
	}
}

func WithAudience(audience ...string) Option {
	return func(c *Config) {
		c.Audience = audience
	}
}

func WithAccessTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.AccessTTL = ttl
	}
}

func WithRefreshTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.RefreshTTL = ttl
	}
}

func WithHMACSecret(secret []byte) Option {
	return func(c *Config) {
		c.HMACSecret = secret
	}
}

func WithStore(store SessionStore) Option {
	return func(c *Config) {
		c.Store = store
	}
}

func WithSQLExec(exec func(ctx context.Context, query string, args ...any) (sql.Result, error)) Option {
	return func(c *Config) {
		c.SQLExec = exec
	}
}

func WithSQLQueryRow(queryRow func(ctx context.Context, query string, args ...any) *sql.Row) Option {
	return func(c *Config) {
		c.SQLQueryRow = queryRow
	}
}
