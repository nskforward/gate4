package retries

import "time"

type config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	MaxJitter    time.Duration
	Backoff      float64
}

func defaultConfig(r *Retry) config {
	cfg := config{
		MaxAttempts:  r.MaxAttempts,
		InitialDelay: r.InitialDelay,
		MaxDelay:     r.MaxDelay,
		MaxJitter:    r.MaxJitter,
		Backoff:      r.Backoff,
	}

	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 10
	}

	if cfg.InitialDelay <= 0 {
		cfg.InitialDelay = 500 * time.Millisecond
	}

	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = 10 * time.Second
	}

	if cfg.MaxJitter <= 0 {
		cfg.MaxJitter = 200 * time.Millisecond
	}

	if cfg.Backoff < 1 {
		cfg.Backoff = 1.5
	}

	return cfg
}
