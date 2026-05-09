package streams

import "iter"

type Stream[T any] struct {
	c           chan T
	unsubscribe func()
	err         *AtomicError
}

func newStream[T any]() *Stream[T] {
	return &Stream[T]{
		c:   make(chan T, 32),
		err: &AtomicError{},
	}
}

func (s *Stream[T]) Close() {
	s.unsubscribe()
	close(s.c)
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

func (s *Stream[T]) registerUnsubscribe(f func()) {
	s.unsubscribe = f
}
