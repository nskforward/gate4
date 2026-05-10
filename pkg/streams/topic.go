package streams

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type topic[T any] struct {
	ctx            context.Context
	cancel         context.CancelFunc
	key            string
	publisher      PublishFunc[T]
	streams        []*Stream[T]
	mx             sync.Mutex
	producerActive atomic.Bool
	attempts       atomic.Int32
	opts           *Opts
	unregister     func(*topic[T])
}

func newTopic[T any](ctx context.Context, opts *Opts, key string, unregister func(*topic[T]), publisher PublishFunc[T]) *topic[T] {
	topicCtx, cancel := context.WithCancel(ctx)
	return &topic[T]{
		ctx:        topicCtx,
		cancel:     cancel,
		key:        key,
		streams:    make([]*Stream[T], 0, 32),
		publisher:  publisher,
		opts:       opts,
		unregister: unregister,
	}
}

func (t *topic[T]) Subscribe(ctx context.Context) *Stream[T] {
	s := newStream(MergedContext(ctx, t.ctx), t.opts.BufferSize, t.unsubscibe)
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

func (t *topic[T]) close() {
	t.unregister(t)
	t.cancel()
}

func (t *topic[T]) startProducer() {
	defer t.producerActive.Store(false)
	for {
		attempts := t.attempts.Add(1)
		err := t.publisher(t.ctx, t.notifyData)
		if err != nil {
			if attempts >= t.opts.MaxRetryAttempts {
				t.notifyErr(err)
				t.close()
				break
			}
			time.Sleep(t.opts.RetryWait(attempts))
			continue
		}
		t.close()
		break
	}
}

func (t *topic[T]) notifyErr(err error) {
	t.mx.Lock()
	defer t.mx.Unlock()
	for _, s := range t.streams {
		if s.ctx.Err() == nil {
			s.err.Store(err)
		}
	}
}

func (t *topic[T]) notifyData(data T) bool {
	if t.attempts.Load() > 0 {
		t.attempts.Store(0)
	}
	t.mx.Lock()
	defer t.mx.Unlock()
	for _, s := range t.streams {
		if s.ctx.Err() == nil {
			s.notify(data)
		}
	}
	return len(t.streams) > 0
}
