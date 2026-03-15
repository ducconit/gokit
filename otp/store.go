package otp

import (
	"context"
	"crypto/sha256"
	"time"
)

type Record struct {
	Purpose     Purpose
	Channel     Channel
	Recipient   string
	CodeSum     [32]byte
	ExpiresAt   time.Time
	Attempts    int
	MaxAttempts int
	LastSentAt  time.Time
	Metadata    map[string]string
}

type Store interface {
	Put(ctx context.Context, r Record) error
	Get(ctx context.Context, purpose Purpose, recipient string) (Record, error)
	Delete(ctx context.Context, purpose Purpose, recipient string) error
}

func sumCode(code string) [32]byte {
	return sha256.Sum256([]byte(code))
}
