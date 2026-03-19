package auth

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
	DB            *sql.DB
	SQLExec       func(ctx context.Context, query string, args ...any) (sql.Result, error)
	SQLQueryRow   func(ctx context.Context, query string, args ...any) *sql.Row
	Dialect       SQLDialect
	SessionsTable string
	BansTable     string
}

func (s *SQLStore) Put(ctx context.Context, sess Session) error {
	if s.SQLExec == nil {
		return fmt.Errorf("auth: sql missing SQLExec")
	}
	st := s.sessionsTable()

	meta, err := json.Marshal(sess.Metadata)
	if err != nil {
		return err
	}

	upd := fmt.Sprintf(
		"update %s set subject_id=%s, subject_type=%s, created_at=%s, expires_at=%s, refresh_sum=%s, metadata=%s where id=%s",
		st,
		s.ph(1), s.ph(2), s.ph(3), s.ph(4), s.ph(5), s.ph(6), s.ph(7),
	)
	res, err := s.SQLExec(ctx, upd,
		sess.SubjectID,
		sess.SubjectType,
		sess.CreatedAt,
		sess.ExpiresAt,
		sess.RefreshTokenSum[:],
		meta,
		sess.ID,
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
		"insert into %s (id, subject_id, subject_type, created_at, expires_at, refresh_sum, metadata) values (%s,%s,%s,%s,%s,%s,%s)",
		st,
		s.ph(1), s.ph(2), s.ph(3), s.ph(4), s.ph(5), s.ph(6), s.ph(7),
	)
	_, err = s.SQLExec(ctx, ins,
		sess.ID,
		sess.SubjectID,
		sess.SubjectType,
		sess.CreatedAt,
		sess.ExpiresAt,
		sess.RefreshTokenSum[:],
		meta,
	)
	return err
}

func (s *SQLStore) Get(ctx context.Context, id string) (Session, error) {
	if s.SQLQueryRow == nil {
		return Session{}, fmt.Errorf("auth: sql missing SQLQueryRow")
	}

	q := fmt.Sprintf(
		"select id, subject_id, subject_type, created_at, expires_at, refresh_sum, metadata from %s where id=%s",
		s.sessionsTable(),
		s.ph(1),
	)
	row := s.SQLQueryRow(ctx, q, id)

	var sess Session
	var sum []byte
	var meta []byte
	if err := row.Scan(&sess.ID, &sess.SubjectID, &sess.SubjectType, &sess.CreatedAt, &sess.ExpiresAt, &sum, &meta); err != nil {
		if err == sql.ErrNoRows {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}
	if len(sum) == 32 {
		copy(sess.RefreshTokenSum[:], sum)
	}
	if len(meta) > 0 {
		_ = json.Unmarshal(meta, &sess.Metadata)
	}
	if time.Now().After(sess.ExpiresAt) {
		return Session{}, ErrSessionNotFound
	}
	return sess, nil
}

func (s *SQLStore) UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error {
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

func (s *SQLStore) Delete(ctx context.Context, id string) error {
	if s.SQLExec == nil {
		return fmt.Errorf("auth: sql missing SQLExec")
	}
	q := fmt.Sprintf("delete from %s where id=%s", s.sessionsTable(), s.ph(1))
	_, err := s.SQLExec(ctx, q, id)
	return err
}

func (s *SQLStore) DeleteBySubject(ctx context.Context, subjectID string) error {
	if s.SQLExec == nil {
		return fmt.Errorf("auth: sql missing SQLExec")
	}
	q := fmt.Sprintf("delete from %s where subject_id=%s", s.sessionsTable(), s.ph(1))
	_, err := s.SQLExec(ctx, q, subjectID)
	return err
}

func (s *SQLStore) BanSubject(ctx context.Context, subjectID string, until time.Time) error {
	if s.SQLExec == nil {
		return fmt.Errorf("auth: sql missing SQLExec")
	}
	bt := s.bansTable()

	upd := fmt.Sprintf("update %s set until_at=%s where subject_id=%s", bt, s.ph(1), s.ph(2))
	res, err := s.SQLExec(ctx, upd, until, subjectID)
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

	ins := fmt.Sprintf("insert into %s (subject_id, until_at) values (%s,%s)", bt, s.ph(1), s.ph(2))
	_, err = s.SQLExec(ctx, ins, subjectID, until)
	return err
}

func (s *SQLStore) IsSubjectBanned(ctx context.Context, subjectID string) (bool, error) {
	if s.SQLQueryRow == nil {
		return false, fmt.Errorf("auth: sql missing SQLQueryRow")
	}

	q := fmt.Sprintf("select until_at from %s where subject_id=%s", s.bansTable(), s.ph(1))
	row := s.SQLQueryRow(ctx, q, subjectID)
	var until time.Time
	if err := row.Scan(&until); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return time.Now().Before(until), nil
}

func (s *SQLStore) ph(n int) string {
	switch s.Dialect {
	case SQLDialectDollar:
		return fmt.Sprintf("$%d", n)
	default:
		return "?"
	}
}

func (s *SQLStore) sessionsTable() string {
	if s.SessionsTable == "" {
		return "auth_sessions"
	}
	return s.SessionsTable
}

func (s *SQLStore) bansTable() string {
	if s.BansTable == "" {
		return "auth_bans"
	}
	return s.BansTable
}
