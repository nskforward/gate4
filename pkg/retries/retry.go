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

func (retry *Retry) Do(ctx context.Context, callback func() error) error {
	attempt := 0
	lastAttempt := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if time.Since(lastAttempt) > retry.cfg.MaxDelay+time.Minute {
			attempt = 1
		} else {
			attempt++
		}

		if retry.cfg.OnAttempt != nil {
			retry.cfg.OnAttempt(attempt)
		}

		err := callback()
		if retry.cfg.OnError != nil {
			retry.cfg.OnError(err)
		}
		if err == nil {
			return nil
		}

		if attempt >= retry.cfg.MaxAttempts {
			return err
		}

		delay := retry.backoff(attempt)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
		}
	}
}

func (retry *Retry) backoff(attempt int) time.Duration {
	delayNano := float64(retry.cfg.InitialDelay) * math.Pow(retry.cfg.BackoffFactor, float64(attempt))
	delay := time.Duration(delayNano) + Jitter(retry.cfg.MaxJitter)
	return min(delay, retry.cfg.MaxDelay)
}

func Jitter(max time.Duration) time.Duration {
	return time.Duration(rand.Int64N(int64(max)))
}
