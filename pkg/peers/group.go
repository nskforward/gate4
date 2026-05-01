package peers

import (
	"slices"
	"sync"
)

type Group[T any] struct {
	key    string
	pubsub *PubSub[T]
	peers  []*Peer[T]
	last   *LastValue[T]
	mx     sync.RWMutex
}

type LastValue[T any] struct {
	exists bool
	value  T
	mx     sync.RWMutex
}

func (last *LastValue[T]) Get() (T, bool) {
	last.mx.RLock()
	defer last.mx.RUnlock()
	return last.value, last.exists
}

func (last *LastValue[T]) Set(value T) {
	last.mx.Lock()
	defer last.mx.Unlock()
	last.value = value
	last.exists = true
}

func NewGroup[T any](key string, pubsub *PubSub[T]) *Group[T] {
	return &Group[T]{
		key:    key,
		pubsub: pubsub,
		peers:  make([]*Peer[T], 0, 8),
		last:   &LastValue[T]{},
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
	v, ok := g.last.Get()
	if ok {
		p.send(v)
	}
	return p
}

func (g *Group[T]) Send(data T) bool {
	g.mx.RLock()
	defer g.mx.RUnlock()
	g.last.Set(data)
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
			g.peers = slices.Delete(g.peers, i, i+1)
			break
		}
	}
}
