package streams

import (
	"context"
	"log/slog"
	"sync"
)

type Store[T any] struct {
	ctx    context.Context
	topics map[string]*topic[T]
	mx     sync.Mutex
	size   int
}

type PublishFunc[T any] func(context.Context, func(data T) bool) error

func NewStore[T any](ctx context.Context, size int) *Store[T] {
	return &Store[T]{
		ctx:    ctx,
		topics: make(map[string]*topic[T]),
		size:   size,
	}
}

func (store *Store[T]) Subscribe(ctx context.Context, key string, publish PublishFunc[T]) *Stream[T] {
	store.mx.Lock()
	topic, ok := store.topics[key]
	if !ok {
		topic = newTopic(store.ctx, key, store.remove, publish)
		store.topics[key] = topic
	}
	store.mx.Unlock()
	if !ok {
		slog.Debug("created a new stream topic", "key", key)
	} else {
		slog.Debug("reuse an existing stream topic", "key", key)
	}
	return topic.Subscribe(ctx, store.size)
}

func (store *Store[T]) remove(t *topic[T]) {
	slog.Debug("delete topic", "key", t.key)
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.topics, t.key)
}
