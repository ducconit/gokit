package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

type Options struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      float64
}

func Do(ctx context.Context, opts Options, fn func(context.Context) error) error {
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = 3
	}
	if opts.BaseDelay <= 0 {
		opts.BaseDelay = 200 * time.Millisecond
	}
	if opts.MaxDelay <= 0 {
		opts.MaxDelay = 5 * time.Second
	}
	if opts.Jitter < 0 {
		opts.Jitter = 0
	}
	if opts.Jitter > 1 {
		opts.Jitter = 1
	}

	var lastErr error
	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			if lastErr != nil {
				return errors.Join(lastErr, err)
			}
			return err
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt == opts.MaxAttempts {
			break
		}

		delay := backoff(attempt, opts.BaseDelay, opts.MaxDelay)
		if opts.Jitter > 0 {
			delay = applyJitter(delay, opts.Jitter)
		}

		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			t.Stop()
			return errors.Join(lastErr, ctx.Err())
		case <-t.C:
		}
	}

	return lastErr
}

func backoff(attempt int, base, max time.Duration) time.Duration {
	pow := math.Pow(2, float64(attempt-1))
	d := time.Duration(float64(base) * pow)
	if d > max {
		return max
	}
	return d
}

func applyJitter(d time.Duration, factor float64) time.Duration {
	if d <= 0 {
		return d
	}
	min := float64(d) * (1 - factor)
	max := float64(d) * (1 + factor)
	n := min + rand.Float64()*(max-min)
	return time.Duration(n)
}
