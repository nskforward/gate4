package streams

import (
	"context"
	"sync"
)

type Store[T any] struct {
	ctx       context.Context
	topics    map[string]*topic[T]
	mx        sync.Mutex
	publisher PublishFunc[T]
}

type PublishFunc[T any] func(ctx context.Context, key string, publish func(data T) bool) error

func NewStore[T any](ctx context.Context, publisher PublishFunc[T]) *Store[T] {
	return &Store[T]{
		ctx:       ctx,
		topics:    make(map[string]*topic[T]),
		publisher: publisher,
	}
}

func (store *Store[T]) Subscribe(ctx context.Context, key string) *Stream[T] {
	store.mx.Lock()
	topic, ok := store.topics[key]
	if !ok {
		topic = newTopic(store.ctx, key, store.remove, store.publisher)
		store.topics[key] = topic
		go store.publishing(topic)
	}
	store.mx.Unlock()
	return topic.Subscribe(ctx, 1)
}

func (store *Store[T]) publishing(t *topic[T]) {
	err := store.publisher(t.ctx, t.key, t.notify)
	if err != nil {
		t.close()
	}
}

func (store *Store[T]) remove(t *topic[T]) {
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.topics, t.key)
}
