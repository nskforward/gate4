package streams

import (
	"context"
	"sync"
)

type Store[T any] struct {
	topics map[string]*topic[T]
	mx     sync.Mutex
	opts   *Opts
}

type PublishFunc[T any] func(context.Context, func(data T) bool) error

func NewStore[T any](opts *Opts) *Store[T] {
	return &Store[T]{
		topics: make(map[string]*topic[T]),
		opts:   initOpts(opts),
	}
}

func (s *Store[T]) Subscribe(ctx context.Context, key string, publisher PublishFunc[T]) *Stream[T] {
	s.mx.Lock()
	topic, ok := s.topics[key]
	if !ok {
		topic = newTopic(s.opts, key, publisher)
		s.topics[key] = topic
	}
	s.mx.Unlock()
	return topic.Subscribe(ctx)
}
