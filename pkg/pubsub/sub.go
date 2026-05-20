package pubsub

import (
	"context"
	"iter"
	"sync"
)

type Sub[T any] struct {
	ctx    context.Context
	once   sync.Once
	closed chan struct{}
	ch     chan T
}

func newSub[T any](ctx context.Context, size int) *Sub[T] {
	sub := &Sub[T]{
		ch:     make(chan T, size),
		closed: make(chan struct{}, 1),
	}
	go func() {
		select {
		case <-ctx.Done():
			sub.Close()
		case <-sub.closed:
			// do nothing
		}
	}()
	return sub
}

func (s *Sub[T]) Close() {
	s.once.Do(func() {
		close(s.closed)
		close(s.ch)
	})
}

// Write returns false if subscriber closed
func (s *Sub[T]) Write(v T) bool {
	for {
		select {
		case <-s.closed:
			return false
		default:
			select {
			case s.ch <- v:
				return true
			default:
				<-s.ch
			}
		}
	}
}

func (s *Sub[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case <-s.closed:
				return
			case v, ok := <-s.ch:
				if !ok {
					return
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}
