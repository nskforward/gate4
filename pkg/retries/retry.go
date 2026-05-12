package retries

import (
	"context"
	"math"
	"math/rand/v2"
	"time"
)

type Retry struct {
	cfg Config
}

func New(cfg Config) *Retry {
	return &Retry{
		cfg: cfg,
	}
}

func (retry *Retry) Do(ctx context.Context, attempt func() error) error {
	attempts := 0
	lastAttempt := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if time.Since(lastAttempt) > retry.cfg.MaxDelay+time.Minute {
			attempts = 1
		} else {
			attempts++
		}

		err := attempt()
		if err == nil {
			return nil
		}

		if attempts >= retry.cfg.MaxAttempts {
			return err
		}

		delay := retry.backoff(attempts)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
		}
	}
}

func (retry *Retry) backoff(attempts int) time.Duration {
	delayNano := float64(retry.cfg.InitialDelay) * math.Pow(retry.cfg.BackoffFactor, float64(attempts))
	delay := time.Duration(delayNano) + time.Duration(rand.Int64N(int64(retry.cfg.MaxJitter)))
	return min(delay, retry.cfg.MaxDelay)
}
