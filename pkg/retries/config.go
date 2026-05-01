package retries

import (
	"math/rand/v2"
	"time"
)

type Config struct {
	InitialDelay  time.Duration // начальная задержка между попытками
	MaxDelay      time.Duration // максимальная задержка (ограничение сверху)
	BackoffFactor float64       // множитель для экспоненциального роста (обычно 1.5–2.0)
	MaxAttempts   int           // максимальное количество попыток (0 = бесконечно)
	JitterFactor  float64       // доля от текущей задержки для добавления случайного джиттера (0.0–1.0)
}

func DefaultConfig() Config {
	return Config{
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxAttempts:   10,
		JitterFactor:  rand.Float64(),
	}
}
