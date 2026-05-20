package peers

import (
	"sync"
)

type PubSub[T any] struct {
	groups map[string]*Group[T]
	mx     sync.Mutex
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		groups: make(map[string]*Group[T]),
	}
}

func (ps *PubSub[T]) LoadOrCreate(key string) (group *Group[T], loaded bool) {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	g, ok := ps.groups[key]
	if ok {
		return g, true
	}
	g = NewGroup(key, ps)
	ps.groups[key] = g
	return g, false
}

func (ps *PubSub[T]) remove(key string) {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	delete(ps.groups, key)
}
