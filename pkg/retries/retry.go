package retries

import (
	"context"
	"math"
	"math/rand/v2"
	"time"
)

type Retry[T any] struct {
	cfg     Config
	connect func() (T, error)
}

func NewRetry[T any](cfg Config, connect func() (T, error)) *Retry[T] {
	return &Retry[T]{
		cfg:     cfg,
		connect: connect,
	}
}

func (retry *Retry[T]) Do(ctx context.Context) (T, error) {
	var attempt int
	var lastErr error

	for {
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		default:
		}

		result, err := retry.connect()
		if err == nil {
			return result, nil
		}
		lastErr = err

		attempt++
		if retry.cfg.MaxAttempts > 0 && attempt >= retry.cfg.MaxAttempts {
			var zero T
			return zero, lastErr
		}

		delayNano := float64(retry.cfg.InitialDelay) * math.Pow(retry.cfg.BackoffFactor, float64(attempt))
		if delayNano > float64(retry.cfg.MaxDelay) {
			delayNano = float64(retry.cfg.MaxDelay)
		}

		delay := time.Duration(delayNano)

		if retry.cfg.JitterFactor > 0 {
			maxJitter := float64(delay) * retry.cfg.JitterFactor
			jitter := rand.Float64() * maxJitter
			delay = time.Duration(float64(delay) + jitter)
			if delay > retry.cfg.MaxDelay {
				delay = retry.cfg.MaxDelay
			}
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			var zero T
			return zero, ctx.Err()
		case <-timer.C:
		}
	}
}
