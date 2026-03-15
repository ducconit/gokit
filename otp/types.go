package otp

import "time"

type Purpose string

type Channel string

type Config struct {
	CodeLength  int
	TTL         time.Duration
	MaxAttempts int
}

type Message struct {
	Purpose   Purpose
	Channel   Channel
	Recipient string
	Code      string
	ExpiresAt time.Time
	Metadata  map[string]string
}
