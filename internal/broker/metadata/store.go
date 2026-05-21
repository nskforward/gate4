package metadata

import (
	"sync"

	"github.com/nskforward/gate4/pkg/types"
)

type Store struct {
	clients map[types.Client]*Metadata
	mx      sync.Mutex
}

func NewStore() *Store {
	return &Store{
		clients: make(map[types.Client]*Metadata),
	}
}

func (store *Store) Get(client types.Client) (*Metadata, error) {
	store.mx.Lock()
	defer store.mx.Unlock()
	metadata, ok := store.clients[client]
	if !ok {
		return store.create(client)
	}
	return metadata, nil
}

func (store *Store) create(client types.Client) (*Metadata, error) {
	metadata, err := NewMetadata(client)
	if err != nil {
		return nil, err
	}
	store.clients[client] = metadata
	return metadata, nil
}
