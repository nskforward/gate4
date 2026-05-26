package retries

import (
	"context"
	"fmt"
	"iter"
	"math/rand"
	"time"
)

type Retry struct {
	// MaxAttempts limits the max number of sequent failures
	// If MaxAttempts = 0 then no limit of attempts
	// Default is 10
	MaxAttempts int

	// InitialDelay is a minimal delay after the first failure
	// Default is 500ms
	InitialDelay time.Duration

	// MaxDelay limits the current delay
	// Default is 10s
	MaxDelay time.Duration

	// MaxJitter a random addition to the current_delay
	// Default is 200ms
	MaxJitter time.Duration

	// Backoff is a multiplication factor: current_delay = previous_delay * Backoff
	// Default is 1.5
	Backoff float64

	err error
}

func (r *Retry) Err() error {
	return r.err
}

func (r *Retry) Range(ctx context.Context) iter.Seq[int] {
	cfg := defaultConfig(r)
	attempt := 0
	delay := cfg.InitialDelay

	return func(yield func(int) bool) {
		for {
			attempt++

			if !yield(attempt) {
				return
			}

			if cfg.MaxAttempts > 0 && attempt >= cfg.MaxAttempts {
				r.err = fmt.Errorf("max attempts (%d) reached", cfg.MaxAttempts)
				return
			}

			select {
			case <-ctx.Done():
				r.err = ctx.Err()
				return

			case <-time.After(delay):
				if delay < cfg.MaxDelay {
					delay = time.Duration(float64(delay) * cfg.Backoff)
					if cfg.MaxJitter > 0 {
						delay = delay + time.Duration(rand.Int63n(int64(cfg.MaxJitter)))
					}
					if delay > cfg.MaxDelay {
						delay = cfg.MaxDelay
					}
				}
			}
		}
	}
}
