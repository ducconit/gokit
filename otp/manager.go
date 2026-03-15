package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

var (
	ErrNotFound    = fmt.Errorf("otp not found")
	ErrExpired     = fmt.Errorf("otp expired")
	ErrMaxAttempts = fmt.Errorf("otp max attempts")
	ErrInvalid     = fmt.Errorf("otp invalid")
)

type Manager struct {
	cfg    Config
	store  Store
	sender Sender
	now    func() time.Time
}

func New(cfg Config, store Store, sender Sender) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("otp: missing store")
	}
	if sender == nil {
		return nil, fmt.Errorf("otp: missing sender")
	}
	if cfg.CodeLength <= 0 {
		cfg.CodeLength = 6
	}
	if cfg.TTL <= 0 {
		cfg.TTL = 5 * time.Minute
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	return &Manager{cfg: cfg, store: store, sender: sender, now: time.Now}, nil
}

func (m *Manager) Request(ctx context.Context, purpose Purpose, channel Channel, recipient string, metadata map[string]string) (Message, error) {
	code, err := generateNumericCode(m.cfg.CodeLength)
	if err != nil {
		return Message{}, err
	}

	now := m.now()
	exp := now.Add(m.cfg.TTL)

	r := Record{
		Purpose:     purpose,
		Channel:     channel,
		Recipient:   recipient,
		CodeSum:     sumCode(code),
		ExpiresAt:   exp,
		Attempts:    0,
		MaxAttempts: m.cfg.MaxAttempts,
		LastSentAt:  now,
		Metadata:    metadata,
	}
	if err := m.store.Put(ctx, r); err != nil {
		return Message{}, err
	}

	msg := Message{
		Purpose:   purpose,
		Channel:   channel,
		Recipient: recipient,
		Code:      code,
		ExpiresAt: exp,
		Metadata:  metadata,
	}
	if err := m.sender.Send(ctx, msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (m *Manager) Verify(ctx context.Context, purpose Purpose, recipient, code string) error {
	r, err := m.store.Get(ctx, purpose, recipient)
	if err != nil {
		return ErrNotFound
	}
	if m.now().After(r.ExpiresAt) {
		_ = m.store.Delete(ctx, purpose, recipient)
		return ErrExpired
	}
	if r.Attempts >= r.MaxAttempts {
		return ErrMaxAttempts
	}

	if r.CodeSum != sumCode(code) {
		r.Attempts++
		_ = m.store.Put(ctx, r)
		return ErrInvalid
	}

	return m.store.Delete(ctx, purpose, recipient)
}

func generateNumericCode(n int) (string, error) {
	max := big.NewInt(10)
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = byte('0' + v.Int64())
	}
	return string(out), nil
}
