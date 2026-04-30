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

func (ps *PubSub[T]) LoadOrCreate(key string) (group *Group[T]) {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	g, ok := ps.groups[key]
	if ok {
		return g
	}
	g = NewGroup(key, ps)
	ps.groups[key] = g
	ps.conf.OnStart(key, g)
	return g
}

func (ps *PubSub[T]) remove(key string) {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	ps.conf.OnStop(key)
	delete(ps.groups, key)
}
