package streams

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
)

type topic[T any] struct {
	ctx            context.Context
	cancel         context.CancelFunc
	key            string
	publish        PublishFunc[T]
	streams        []*Stream[T]
	mx             sync.Mutex
	producerActive atomic.Bool
	unregister     func(*topic[T])
	last           *T
	once           sync.Once
}

func newTopic[T any](ctx context.Context, key string, unregister func(*topic[T]), publish PublishFunc[T]) *topic[T] {
	topicCtx, cancel := context.WithCancel(ctx)
	return &topic[T]{
		ctx:        topicCtx,
		cancel:     cancel,
		key:        key,
		streams:    make([]*Stream[T], 0, 32),
		publish:    publish,
		unregister: unregister,
	}
}

func (t *topic[T]) Subscribe(ctx context.Context, size int) *Stream[T] {
	s := newStream(MergedContext(ctx, t.ctx), size, t.unsubscibe)
	t.mx.Lock()
	t.streams = append(t.streams, s)
	t.mx.Unlock()
	if t.producerActive.CompareAndSwap(false, true) {
		go t.startProducer()
	}
	if t.last != nil {
		s.notify(*t.last)
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
	if len(t.streams) == 0 {
		t.close()
	}
}

func (t *topic[T]) close() {
	t.once.Do(func() {
		t.unregister(t)
		t.cancel()
	})
}

func (t *topic[T]) startProducer() {
	defer t.producerActive.Store(false)
	err := t.publish(t.ctx, t.key, t.notify)
	if err != nil {
		t.notifyErr(err)
	}
	t.close()
}

func (t *topic[T]) notify(data T) bool {
	t.mx.Lock()
	defer t.mx.Unlock()
	for _, s := range t.streams {
		if s.ctx.Err() == nil {
			s.notify(data)
		}
	}
	t.last = &data
	return len(t.streams) > 0
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
