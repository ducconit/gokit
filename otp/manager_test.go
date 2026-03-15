package otp

import (
	"context"
	"testing"
	"time"
)

type dummySender struct {
	last Message
}

func (d *dummySender) Send(ctx context.Context, msg Message) error {
	_ = ctx
	d.last = msg
	return nil
}

func TestOTPRequestVerify(t *testing.T) {
	ctx := context.Background()

	store := NewMemoryStore()
	sender := &dummySender{}
	m, err := New(Config{
		CodeLength:  6,
		TTL:         1 * time.Minute,
		MaxAttempts: 2,
	}, store, sender)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	msg, err := m.Request(ctx, "forgot_password", "sms", "+840000", nil)
	if err != nil {
		t.Fatalf("Request: %v", err)
	}
	if msg.Code == "" {
		t.Fatalf("expected code")
	}

	if err := m.Verify(ctx, "forgot_password", "+840000", "000000"); err != ErrInvalid {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}

	if err := m.Verify(ctx, "forgot_password", "+840000", msg.Code); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if err := m.Verify(ctx, "forgot_password", "+840000", msg.Code); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
