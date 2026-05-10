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
	streams        []*Stream[T]
	mx             sync.Mutex
	producerActive atomic.Bool
	attempts       atomic.Int32
	opts           *Opts
}

func newTopic[T any](opts *Opts, key string, publisher PublishFunc[T]) *topic[T] {
	return &topic[T]{
		key:       key,
		streams:   make([]*Stream[T], 0, 32),
		publisher: publisher,
		opts:      opts,
	}
}

func (t *topic[T]) Subscribe(ctx context.Context) *Stream[T] {
	s := newStream(ctx, t.opts.BufferSize, t.unsubscibe)
	t.mx.Lock()
	t.streams = append(t.streams, s)
	t.mx.Unlock()
	if t.producerActive.CompareAndSwap(false, true) {
		go t.startProducer()
	}
	return s
}

func (t *topic[T]) unsubscibe(stream *Stream[T]) {
	t.mx.Lock()
	defer t.mx.Unlock()
	for i, s := range t.streams {
		if s == stream {
			t.streams = slices.Delete(t.streams, i, i+1)
			break
		}
	}
}

func (t *topic[T]) close(err error) {
	t.mx.Lock()
	streams := make([]*Stream[T], 0, len(t.streams))
	for _, s := range t.streams {
		streams = append(streams, s)
	}
	t.mx.Unlock()
	for _, s := range streams {
		if err != nil {
			s.err.Store(err)
		}
		s.Close()
	}
}

func (t *topic[T]) startProducer() {
	defer t.producerActive.Store(false)
	for {
		attempts := t.attempts.Add(1)
		err := t.publisher(context.Background(), t.notify)
		if err != nil {
			if attempts >= t.opts.MaxRetryAttempts {
				t.close(err)
				break
			}
			time.Sleep(t.opts.RetryWait(attempts))
			continue
		}
		t.close(nil)
		break
	}
}

func (t *topic[T]) notify(data T) bool {
	t.attempts.Store(0)
	t.mx.Lock()
	defer t.mx.Unlock()
	for _, s := range t.streams {
		s.notify(data)
	}
	return len(t.streams) > 0
}
