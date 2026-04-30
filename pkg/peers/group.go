package peers

import (
	"slices"
	"sync"
)

type Group[T any] struct {
	key    string
	pubsub *PubSub[T]
	peers  []*Peer[T]
	mx     sync.RWMutex
}

func NewGroup[T any](key string, pubsub *PubSub[T]) *Group[T] {
	return &Group[T]{
		key:    key,
		pubsub: pubsub,
		peers:  make([]*Peer[T], 0, 8),
	}
}

func (g *Group[T]) Close() {
	g.mx.Lock()
	defer g.mx.Unlock()
	if g.pubsub != nil {
		g.pubsub.remove(g.key)
	}
	for _, p := range g.peers {
		p.Close()
	}
}

func (g *Group[T]) NewPeer() *Peer[T] {
	g.mx.Lock()
	defer g.mx.Unlock()
	p := NewPeer(8, g)
	g.peers = append(g.peers, p)
	return p
}

func (g *Group[T]) Send(data T) bool {
	g.mx.RLock()
	defer g.mx.RUnlock()
	notified := 0
	for _, p := range g.peers {
		p.send(data)
		notified++
	}
	return notified > 0
}

func (g *Group[T]) remove(p *Peer[T]) {
	g.mx.Lock()
	defer g.mx.Unlock()
	for i, p1 := range g.peers {
		if p1 == p {
			g.peers = slices.Delete(g.peers, i, 1)
			break
		}
	}
}
