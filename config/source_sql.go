package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type SQLSource struct {
	DB        *sql.DB
	Query     string // Query for Load (e.g. "SELECT value FROM configs WHERE key = ?")
	SaveQuery string // Query for Save (e.g. "UPDATE configs SET value = ? WHERE key = ?")
	Args      []any

	// ScanAll when true, Query should return multiple rows of (key, value) pairs.
	// If Query is empty and ScanAll is true, it will use default "SELECT key, value FROM configs".
	ScanAll bool
}

func (s *SQLSource) Load(ctx context.Context) ([]byte, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("config: sql missing db")
	}

	if s.ScanAll {
		query := s.Query
		if query == "" {
			query = "SELECT `key`, `value` FROM `configs`"
		}

		rows, err := s.DB.QueryContext(ctx, query, s.Args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		configs := make(map[string]any)
		for rows.Next() {
			var key string
			var val []byte
			if err := rows.Scan(&key, &val); err != nil {
				return nil, err
			}

			// Try to unmarshal value if it looks like JSON, otherwise keep as string
			var jsonVal any
			if err := json.Unmarshal(val, &jsonVal); err == nil {
				configs[key] = jsonVal
			} else {
				configs[key] = string(val)
			}
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return json.Marshal(configs)
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

func (s *SQLSource) Write(ctx context.Context, b []byte) error {
	if s.DB == nil {
		return fmt.Errorf("config: sql missing db")
	}

	if s.ScanAll {
		// When ScanAll is true, we expect b to be a JSON map of configs.
		// We'll upsert each key-value pair.
		var configs map[string]any
		if err := json.Unmarshal(b, &configs); err != nil {
			return fmt.Errorf("config: failed to unmarshal for SQL write: %w", err)
		}

		tx, err := s.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		query := s.SaveQuery
		if query == "" {
			// Assuming a standard table named configs with key and value columns
			query = "INSERT INTO `configs` (`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)"
		}

		for k, v := range configs {
			vb, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, query, k, vb); err != nil {
				return err
			}
		}

		return tx.Commit()
	}

	if s.SaveQuery == "" {
		return fmt.Errorf("config: sql missing save query")
	}

	_, err := s.DB.ExecContext(ctx, s.SaveQuery, append(s.Args, b)...)
	return err
}
