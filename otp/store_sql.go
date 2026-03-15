package otp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type SQLDialect int

const (
	SQLDialectQuestion SQLDialect = iota
	SQLDialectDollar
)

type SQLStore struct {
	DB      *sql.DB
	Dialect SQLDialect
	Table   string
}

func NewSQLStore(db *sql.DB, dialect SQLDialect, table string) *SQLStore {
	if table == "" {
		table = "otp_records"
	}
	return &SQLStore{
		DB:      db,
		Dialect: dialect,
		Table:   table,
	}
}

func (s *SQLStore) Put(ctx context.Context, r Record) error {
	if s.DB == nil {
		return fmt.Errorf("otp: sql missing db")
	}

	meta, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}

	upd := fmt.Sprintf(
		"update %s set channel=%s, code_sum=%s, expires_at=%s, attempts=%s, max_attempts=%s, last_sent_at=%s, metadata=%s where purpose=%s and recipient=%s",
		s.Table,
		s.ph(1), s.ph(2), s.ph(3), s.ph(4), s.ph(5), s.ph(6), s.ph(7), s.ph(8), s.ph(9),
	)
	res, err := s.DB.ExecContext(ctx, upd,
		r.Channel,
		r.CodeSum[:],
		r.ExpiresAt,
		r.Attempts,
		r.MaxAttempts,
		r.LastSentAt,
		meta,
		r.Purpose,
		r.Recipient,
	)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		ra = 0
	}
	if ra > 0 {
		return nil
	}

	ins := fmt.Sprintf(
		"insert into %s (purpose, channel, recipient, code_sum, expires_at, attempts, max_attempts, last_sent_at, metadata) values (%s,%s,%s,%s,%s,%s,%s,%s,%s)",
		s.Table,
		s.ph(1), s.ph(2), s.ph(3), s.ph(4), s.ph(5), s.ph(6), s.ph(7), s.ph(8), s.ph(9),
	)
	_, err = s.DB.ExecContext(ctx, ins,
		r.Purpose,
		r.Channel,
		r.Recipient,
		r.CodeSum[:],
		r.ExpiresAt,
		r.Attempts,
		r.MaxAttempts,
		r.LastSentAt,
		meta,
	)
	return err
}

func (s *SQLStore) Get(ctx context.Context, purpose Purpose, recipient string) (Record, error) {
	if s.DB == nil {
		return Record{}, fmt.Errorf("otp: sql missing db")
	}

	q := fmt.Sprintf(
		"select purpose, channel, recipient, code_sum, expires_at, attempts, max_attempts, last_sent_at, metadata from %s where purpose=%s and recipient=%s",
		s.Table,
		s.ph(1), s.ph(2),
	)
	row := s.DB.QueryRowContext(ctx, q, purpose, recipient)

	var r Record
	var sum []byte
	var meta []byte
	if err := row.Scan(&r.Purpose, &r.Channel, &r.Recipient, &sum, &r.ExpiresAt, &r.Attempts, &r.MaxAttempts, &r.LastSentAt, &meta); err != nil {
		if err == sql.ErrNoRows {
			return Record{}, ErrNotFound
		}
		return Record{}, err
	}
	if len(sum) == 32 {
		copy(r.CodeSum[:], sum)
	}
	if len(meta) > 0 {
		_ = json.Unmarshal(meta, &r.Metadata)
	}
	if time.Now().After(r.ExpiresAt) {
		return Record{}, ErrExpired
	}
	return r, nil
}

func (s *SQLStore) Delete(ctx context.Context, purpose Purpose, recipient string) error {
	if s.DB == nil {
		return fmt.Errorf("otp: sql missing db")
	}
	q := fmt.Sprintf("delete from %s where purpose=%s and recipient=%s", s.Table, s.ph(1), s.ph(2))
	_, err := s.DB.ExecContext(ctx, q, purpose, recipient)
	return err
}

func (s *SQLStore) ph(i int) string {
	if s.Dialect == SQLDialectDollar {
		return fmt.Sprintf("$%d", i)
	}
	return "?"
}
