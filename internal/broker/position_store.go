package broker

import (
	"fmt"
	"sync"

	"github.com/nskforward/gate4/pkg/types"
)

type PositionStore struct {
	accounts map[Client][]types.Position
	mx       sync.Mutex
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		accounts: make(map[Client][]types.Position),
	}
}

func (store *PositionStore) Get(client Client) ([]types.Position, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	positions, ok := store.accounts[client]
	if ok {
		return positions, nil
	}

	return nil, fmt.Errorf("not implemented")
}
