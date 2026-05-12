package retries

import (
	"time"
)

type Config struct {
	InitialDelay  time.Duration // начальная задержка между попытками
	MaxDelay      time.Duration // максимальная задержка (ограничение сверху)
	BackoffFactor float64       // множитель для экспоненциального роста (обычно 1.5–2.0)
	MaxAttempts   int           // максимальное количество попыток (0 = бесконечно)
	MaxJitter     time.Duration // random addition to backoff
	OnBefore      func(int)
	OnAfter       func(error)
}

func DefaultConfig() Config {
	return Config{
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxAttempts:   10,
		MaxJitter:     time.Second,
	}
}
