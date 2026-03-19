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
}
