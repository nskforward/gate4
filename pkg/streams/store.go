package streams

import (
	"context"
	"sync"
)

type Store[T any] struct {
	ctx    context.Context
	topics map[string]*topic[T]
	mx     sync.Mutex
	opts   *Opts
}

type PublishFunc[T any] func(context.Context, func(data T) bool) error

func NewStore[T any](ctx context.Context, opts *Opts) *Store[T] {
	return &Store[T]{
		ctx:    ctx,
		topics: make(map[string]*topic[T]),
		opts:   initOpts(opts),
	}
}

func (store *Store[T]) Subscribe(ctx context.Context, key string, publisher PublishFunc[T]) *Stream[T] {
	store.mx.Lock()
	topic, ok := store.topics[key]
	if !ok {
		topic = newTopic(store.ctx, store.opts, key, store.remove, publisher)
		store.topics[key] = topic
	}
	store.mx.Unlock()
	return topic.Subscribe(ctx)
}

func (store *Store[T]) remove(t *topic[T]) {
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.topics, t.key)
}
