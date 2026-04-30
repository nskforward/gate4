package peers

import "context"

type Peer[T any] struct {
	group   *Group[T]
	channel chan T
}

func NewPeer[T any](g *Group[T], size int) *Peer[T] {
	return &Peer[T]{
		group:   g,
		channel: make(chan T, size),
	}
}

func (p *Peer[T]) Close() {
	p.group.remove(p)
	close(p.channel)
}

func (p *Peer[T]) Wait(ctx context.Context) (T, bool) {
	select {
	case <-ctx.Done():
		var def T
		return def, false
	case data, ok := <-p.channel:
		return data, ok
	}
}

func (p *Peer[T]) send(data T) {
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
