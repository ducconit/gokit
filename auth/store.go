package auth

import (
	"context"
	"time"
)

type Session struct {
	ID              string            `json:"id"`
	SubjectID       string            `json:"sid"`
	SubjectType     string            `json:"stp,omitempty"`
	CreatedAt       time.Time         `json:"cat"`
	ExpiresAt       time.Time         `json:"exp"`
	RefreshTokenSum [32]byte          `json:"rfs"`
	Metadata        map[string]string `json:"meta,omitempty"`
}

type SessionStore interface {
	Put(ctx context.Context, s Session) error
	Get(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
	DeleteBySubject(ctx context.Context, subjectID string) error
	UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error

	BanSubject(ctx context.Context, subjectID string, until time.Time) error
	IsSubjectBanned(ctx context.Context, subjectID string) (bool, error)
}
