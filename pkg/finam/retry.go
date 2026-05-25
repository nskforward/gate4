package finam

import (
	"context"
	"time"
)

type Retry[T any] struct {
	retryFunc RetryFunc[T]
	OnSuccess func(attempt int)
	OnFailure func(err error, attempt int)
}

type RetryFunc[T any] func() (T, error)

func NewRetry[T any](retryFunc RetryFunc[T]) *Retry[T] {
	return &Retry[T]{
		retryFunc: retryFunc,
	}
}

func (r *Retry[T]) Do(ctx context.Context) (T, error) {
	attempts := 0
	sleep := 100 * time.Millisecond

	for {
		attempts++

		result, err := r.retryFunc()
		if err == nil {
			if r.OnSuccess != nil {
				r.OnSuccess(attempts)
			}
			return result, nil
		}

		if r.OnFailure != nil {
			r.OnFailure(err, attempts)
		}

		if attempts > 10 {
			return result, err
		}

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(sleep):
			sleep = sleep * 2
		}
	}
}
