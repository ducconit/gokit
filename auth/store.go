package auth

import (
	"context"
	"time"
)

type Session struct {
	ID              string            `json:"id"`
	UserID          string            `json:"uid"`
	UserType        string            `json:"utp,omitempty"`
	CreatedAt       time.Time         `json:"cat"`
	ExpiresAt       time.Time         `json:"exp"`
	RefreshTokenSum [32]byte          `json:"rfs"`
	RevokedAt       *time.Time        `json:"rvk,omitempty"`
	Metadata        map[string]string `json:"meta,omitempty"`
}

type SessionStore interface {
	Put(ctx context.Context, s Session) error
	Get(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
	DeleteByUser(ctx context.Context, userID string) error
	UpdateMetadata(ctx context.Context, id string, metadata map[string]string) error

	BanUser(ctx context.Context, userID string, until time.Time) error
	IsUserBanned(ctx context.Context, userID string) (bool, error)
}
