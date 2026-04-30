package peers

import (
	"slices"
	"sync"
)

type Group[T any] struct {
	peers []*Peer[T]
	mx    sync.RWMutex
}

func NewGroup[T any]() *Group[T] {
	return &Group[T]{
		peers: make([]*Peer[T], 0, 8),
	}
}

func (g *Group[T]) Close() {
	g.mx.Lock()
	defer g.mx.Unlock()
	for _, p := range g.peers {
		p.Close()
	}
}

func (g *Group[T]) NewPeer() (*Peer[T], int) {
	g.mx.Lock()
	defer g.mx.Unlock()
	p := NewPeer(g, 8)
	g.peers = append(g.peers, p)
	return p, len(g.peers) - 1
}

func (g *Group[T]) Send(data T) {
	g.mx.RLock()
	defer g.mx.RUnlock()
	for _, p := range g.peers {
		p.send(data)
	}
}

func (g *Group[T]) remove(p *Peer[T]) int {
	g.mx.Lock()
	defer g.mx.Unlock()
	for i, p1 := range g.peers {
		if p1 == p {
			g.peers = slices.Delete(g.peers, i, 1)
			return len(g.peers)
		}
	}
	return -1
}
