package streams

import (
	"context"
	"sync"
)

type Store[T any] struct {
	topics map[string]*topic[T]
	mx     sync.Mutex
}

type Publish[T any] func(context.Context, func(data T) bool) error

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		topics: make(map[string]*topic[T]),
	}
}

func (s *Store[T]) Subscribe(key string, publish Publish[T]) *Stream[T] {
	s.mx.Lock()
	topic, ok := s.topics[key]
	if !ok {
		topic = newTopic(key, publish)
		s.topics[key] = topic
	}
	s.mx.Unlock()
	return topic.Subscribe()
}
