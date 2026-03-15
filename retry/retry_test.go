package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoRetries(t *testing.T) {
	ctx := context.Background()

	var n int
	err := Do(ctx, Options{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    2 * time.Millisecond,
		Jitter:      0,
	}, func(ctx context.Context) error {
		_ = ctx
		n++
		if n < 3 {
			return errors.New("fail")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 attempts, got %d", n)
	}
}

func TestDoContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Do(ctx, Options{MaxAttempts: 3}, func(ctx context.Context) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
