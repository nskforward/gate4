package broker

import (
	"sync"
)

type clientStore struct {
	clients map[Client]*metadata
	mx      sync.Mutex
}

type metadata struct {
	positions     *positionStore
	quotes        *quoteStore
	accountTrades *acountTradeStore
	schedules     *scheduleStore
}

func newClientStore() *clientStore {
	return &clientStore{
		clients: make(map[Client]*metadata),
	}
}

func (store *clientStore) Get(client Client) (*metadata, error) {
	store.mx.Lock()
	defer store.mx.Unlock()
	metadata, ok := store.clients[client]
	if !ok {
		return store.create(client)
	}
	return metadata, nil
}

func (store *clientStore) create(client Client) (*metadata, error) {
	positions, err := newPositionStore(client)
	if err != nil {
		return nil, err
	}
	metadata := &metadata{
		positions: positions,
		quotes:    newQuoteStore(client),
		schedules: newScheduleStore(client),
	}
	store.clients[client] = metadata
	return metadata, nil
}
