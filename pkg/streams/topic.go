package streams

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type topic[T any] struct {
	key            string
	publishing     Publish[T]
	subscribers    []*Stream[T]
	mx             sync.Mutex
	producerActive bool
	attempts       atomic.Int32
}

func newTopic[T any](key string, publishFunc Publish[T]) *topic[T] {
	return &topic[T]{
		key:         key,
		subscribers: make([]*Stream[T], 0, 32),
		publishing:  publishFunc,
	}
}

func (t *topic[T]) Subscribe() *Stream[T] {
	stream := newStream[T]()

	t.mx.Lock()
	t.subscribers = append(t.subscribers, stream)
	if !t.producerActive {
		t.producerActive = true
		go t.startProducer()
	}
	t.mx.Unlock()

	stream.registerUnsubscribe(func() {
		t.unsubscibe(stream)
	})

	return stream
}

func (t *topic[T]) unsubscibe(subscriber *Stream[T]) {
	t.mx.Lock()
	defer t.mx.Unlock()
	for i, sub := range t.subscribers {
		if sub == subscriber {
			t.subscribers = slices.Delete(t.subscribers, i, i+1)
			break
		}
	}
}

func (t *topic[T]) fail(err error) {
	t.mx.Lock()
	subscribers := make([]*Stream[T], 0, len(t.subscribers))
	for _, s := range t.subscribers {
		subscribers = append(subscribers, s)
	}
	t.mx.Unlock()

	for _, s := range subscribers {
		s.err.Store(err)
		s.Close()
	}
}

func (t *topic[T]) startProducer() {
	defer func() {
		t.mx.Lock()
		t.producerActive = false
		t.mx.Unlock()
	}()

	for {
		attempts := t.attempts.Add(1)
		err := t.publishing(context.Background(), t.notify)
		if err != nil {
			slog.Error("stream producer error", "topic", t.key, "error", err.Error(), "attempts", attempts)
			if attempts > 5 {
				slog.Warn("close stream topic due to max attempts reaching", "topic", t.key, "attempts", attempts)
				t.fail(err)
				break
			}
			time.Sleep(time.Duration(attempts) * time.Second)
			continue
		}
		slog.Info("stream producer stopped due to no subscribers", "topic", t.key)
		break
	}
}

func (t *topic[T]) notify(data T) bool {
	t.attempts.Store(0)
	t.mx.Lock()
	defer t.mx.Unlock()

	for _, subscriber := range t.subscribers {
		subscriber.notify(data)
	}

	return len(t.subscribers) > 0
}
