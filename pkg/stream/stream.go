package stream

import (
	"context"
	"iter"
	"sync"
)

type Stream[T any] struct {
	ctx       context.Context
	cancel    context.CancelFunc
	items     chan T
	err       AtomicError
	closeOnce sync.Once
}

type SourceFunc[T any] func() ([]T, error)

func NewStream[T any](ctx context.Context, size int, getter SourceFunc[T]) *Stream[T] {
	streamCtx, cancel := context.WithCancel(ctx)
	s := &Stream[T]{
		ctx:    streamCtx,
		cancel: cancel,
		items:  make(chan T),
	}
	go s.watch(getter)
	return s
}

func (s *Stream[T]) Close() {
	s.closeOnce.Do(s.cancel)
}

func (s *Stream[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case <-s.ctx.Done():
				return
			case item, ok := <-s.items:
				if !ok {
					return
				}
				if !yield(item) {
					return
				}
			}
		}
	}
}

func (s *Stream[T]) Err() error {
	return s.err.Load()
}

func (s *Stream[T]) watch(getter SourceFunc[T]) {
	defer close(s.items)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			items, err := getter()
			if err != nil {
				s.err.Store(err)
				return
			}
			for _, item := range items {
				select {
				case <-s.ctx.Done():
					return
				case s.items <- item:
				}
			}
		}
	}
}
