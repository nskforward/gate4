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

func (s *Store[T]) Subscribe(ctx context.Context, key string, publisher PublishFunc[T]) *Stream[T] {
	s.mx.Lock()
	topic, ok := s.topics[key]
	if !ok {
		topic = newTopic(s.ctx, s.opts, key, s.remove, publisher)
		s.topics[key] = topic
	}
	s.mx.Unlock()
	return topic.Subscribe(ctx)
}

func (s *Store[T]) remove(t *topic[T]) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.topics, t.key)
}
