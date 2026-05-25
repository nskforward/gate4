package retry

import (
	"context"
	"iter"
	"math/rand"
	"time"
)

type Retry struct {
	// MaxAttempts limits the max number of sequent failures
	// Default is 10 if not specified
	MaxAttempts int

	// InitialDelay is a minimal delay after the first failure
	// Default is 100ms if not specified
	InitialDelay time.Duration

	// MaxDelay limits the current delay
	// Default is 30s if not specified
	MaxDelay time.Duration

	// MaxJitter a random addition to the current_delay
	// Default is 0 (no jitter) if not specified
	MaxJitter time.Duration

	// Backoff is a multiplication factor: current_delay = previous_delay * Backoff
	// Default is 1.5 if not specified
	Backoff float64

	err            error
	currentAttempt int
	currentDelay   time.Duration
}

func (r *Retry) Err() error {
	return r.err
}

func (r *Retry) Range(ctx context.Context) iter.Seq[Attempt] {
	cfg := defaultConfig(r)

	r.currentAttempt = 0
	r.currentDelay = cfg.InitialDelay

	return func(yield func(Attempt) bool) {
		for {
			r.currentAttempt++

			if !yield(Attempt{r, &cfg}) {
				return
			}

			if r.currentAttempt >= cfg.MaxAttempts {
				return
			}

			select {
			case <-ctx.Done():
				r.err = ctx.Err()
				return

			case <-time.After(r.currentDelay):
				if r.currentDelay < cfg.MaxDelay {
					r.currentDelay = time.Duration(float64(r.currentDelay) * cfg.Backoff)
					if cfg.MaxJitter > 0 {
						r.currentDelay = r.currentDelay + time.Duration(rand.Int63n(int64(cfg.MaxJitter)))
					}
					if r.currentDelay > cfg.MaxDelay {
						r.currentDelay = cfg.MaxDelay
					}
				}
			}
		}
	}
}
