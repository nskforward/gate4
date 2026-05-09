package streams

import (
	"context"
	"iter"
	"sync"
)

type Stream[T any] struct {
	ctx         context.Context
	c           chan T
	unsubscribe func(*Stream[T])
	err         *AtomicError
	once        sync.Once
}

func (s *Stream[T]) Close() {
	s.once.Do(func() {
		s.unsubscribe(s)
		close(s.c)
	})
}

func (s *Stream[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case <-s.ctx.Done():
				s.err.Store(s.ctx.Err())
				return
			case data, ok := <-s.c:
				if !ok {
					return
				}
				if !yield(data) {
					return
				}
			}
		}
	}
}

func (s *Stream[T]) Err() error {
	return s.err.Load()
}

func newStream[T any](ctx context.Context, unsubscribe func(*Stream[T])) *Stream[T] {
	return &Stream[T]{
		ctx:         ctx,
		c:           make(chan T, 16),
		err:         &AtomicError{},
		unsubscribe: unsubscribe,
	}
}

func (s *Stream[T]) notify(data T) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case s.c <- data:
			return
		default:
			select {
			case <-s.ctx.Done():
				return
			case <-s.c:
			default:
			}
		}
	}
}
