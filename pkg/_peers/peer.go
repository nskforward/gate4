package peers

import (
	"context"
	"sync"
)

type Peer[T any] struct {
	group   *Group[T]
	channel chan T
	mx      sync.Mutex
}

func NewPeer[T any](size int, g *Group[T]) *Peer[T] {
	return &Peer[T]{
		group:   g,
		channel: make(chan T, size),
	}
}

func (p *Peer[T]) Close() {
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.group != nil {
		p.group.remove(p)
	}

	close(p.channel)
}

func (p *Peer[T]) Read(ctx context.Context) (T, bool) {
	select {
	case <-ctx.Done():
		var def T
		return def, false
	case data, ok := <-p.channel:
		return data, ok
	}
}

func (p *Peer[T]) send(data T) {
	p.mx.Lock()
	defer p.mx.Unlock()

	for {
		select {
		case p.channel <- data:
			return
		default:
			select {
			case <-p.channel:
			default:
			}
		}
	}
}

func (p *Peer[T]) close() {
	p.mx.Lock()
	defer p.mx.Unlock()
	close(p.channel)
}
