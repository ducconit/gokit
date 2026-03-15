package config

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLSource struct {
	DB    *sql.DB
	Query string
	Args  []any
}

func (s *SQLSource) Load(ctx context.Context) ([]byte, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("config: sql missing db")
	}
	if s.Query == "" {
		return nil, fmt.Errorf("config: sql missing query")
	}

	row := s.DB.QueryRowContext(ctx, s.Query, s.Args...)
	var b []byte
	if err := row.Scan(&b); err != nil {
		return nil, err
	}
	return append([]byte(nil), b...), nil
}
