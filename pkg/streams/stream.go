package streams

import (
	"iter"
	"sync"
)

type Stream[T any] struct {
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
		for data := range s.c {
			if !yield(data) {
				break
			}
		}
	}
}

func (s *Stream[T]) Err() error {
	return s.err.Load()
}

func newStream[T any](unsubscribe func(*Stream[T])) *Stream[T] {
	return &Stream[T]{
		c:           make(chan T, 16),
		err:         &AtomicError{},
		unsubscribe: unsubscribe,
	}
}

func (s *Stream[T]) notify(data T) {
	for {
		select {
		case s.c <- data:
			return
		default:
			select {
			case <-s.c:
			default:
			}
		}
	}
}
