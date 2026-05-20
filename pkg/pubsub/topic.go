package pubsub

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
)

type Topic[T any] struct {
	closed atomic.Bool
	size   int
	subs   []*Sub[T]
	mx     sync.Mutex
	last   atomic.Pointer[T]
}

func NewTopic[T any](size int) *Topic[T] {
	if size < 1 {
		panic("topic size must be a positive number and greater than 0")
	}
	return &Topic[T]{
		size: size,
		subs: make([]*Sub[T], 0, 16),
	}
}

func (t *Topic[T]) Close() {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.closed.Store(true)
	for _, sub := range t.subs {
		sub.Close()
	}
	t.subs = t.subs[:0]
}

func (t *Topic[T]) Last() *T {
	return t.last.Load()
}

func (t *Topic[T]) Publish(v T) int {
	if t.closed.Load() {
		return 0
	}
	t.last.Store(&v)
	subs := t.copy()
	i := 0
	for _, sub := range subs {
		if !sub.Write(v) {
			t.delete(sub)
		} else {
			i++
		}
	}
	return i
}

func (t *Topic[T]) Subscribe(ctx context.Context) (*Sub[T], error) {
	if t.closed.Load() {
		return nil, fmt.Errorf("cannot subscribe on closed topic")
	}
	sub := newSub[T](ctx, t.size)
	t.mx.Lock()
	t.subs = append(t.subs, sub)
	t.mx.Unlock()
	return sub, nil
}

func (t *Topic[T]) delete(v *Sub[T]) {
	t.mx.Lock()
	defer t.mx.Unlock()

	for i, sub := range t.subs {
		if sub == v {
			t.subs = slices.Delete(t.subs, i, i+1)
			return
		}
	}
}

func (t *Topic[T]) copy() []*Sub[T] {
	t.mx.Lock()
	defer t.mx.Unlock()

	result := make([]*Sub[T], 0, len(t.subs))
	for _, sub := range t.subs {
		result = append(result, sub)
	}

	return result
}
