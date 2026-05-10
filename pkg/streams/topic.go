package streams

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type topic[T any] struct {
	key            string
	publisher      PublishFunc[T]
	subscribers    []*Stream[T]
	mx             sync.Mutex
	producerActive atomic.Bool
	attempts       atomic.Int32
	opts           *Opts
}

func newTopic[T any](opts *Opts, key string, publisher PublishFunc[T]) *topic[T] {
	return &topic[T]{
		key:         key,
		subscribers: make([]*Stream[T], 0, 32),
		publisher:   publisher,
		opts:        opts,
	}
}

func (t *topic[T]) Subscribe(ctx context.Context) *Stream[T] {
	stream := newStream(ctx, t.opts.BufferSize, t.unsubscibe)
	t.mx.Lock()
	t.subscribers = append(t.subscribers, stream)
	t.mx.Unlock()
	if t.producerActive.CompareAndSwap(false, true) {
		go t.startProducer()
	}
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

func (t *topic[T]) close(err error) {
	t.mx.Lock()
	subscribers := make([]*Stream[T], 0, len(t.subscribers))
	for _, s := range t.subscribers {
		subscribers = append(subscribers, s)
	}
	t.mx.Unlock()

	for _, s := range subscribers {
		if err != nil {
			s.err.Store(err)
		}
		s.Close()
	}
}

func (t *topic[T]) startProducer() {
	defer t.producerActive.Store(false)
	//slog.Debug("started stream producer", "key", t.key)

	for {
		attempts := t.attempts.Add(1)

		err := t.publisher(context.Background(), t.notify)
		if err != nil {
			// slog.Error("stream producer error", "topic", t.key, "error", err.Error(), "attempts", attempts)
			if attempts >= t.opts.MaxRetryAttempts {
				t.close(err)
				break
			}
			time.Sleep(t.opts.RetryWait(attempts))
			continue
		}
		//slog.Debug("stream producer stopped due to no subscribers", "topic", t.key)
		t.close(nil)
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
