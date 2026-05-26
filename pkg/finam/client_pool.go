package finam

import (
	"sync"
)

type ClientPool struct {
	clients map[string]*Client
	mx      sync.Mutex
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		clients: make(map[string]*Client),
	}
}

func (pool *ClientPool) DeleteClient(accountID string) error {
	pool.mx.Lock()
	defer pool.mx.Unlock()

	delete(pool.clients, accountID)
	return nil
}

func (pool *ClientPool) GetOrCreateClient(creds *Creds) (*Client, error) {
	pool.mx.Lock()
	defer pool.mx.Unlock()

	client, ok := pool.clients[creds.AccountID]
	if ok {
		return client, nil
	}

	client, err := NewClient(creds)
	if err != nil {
		return nil, err
	}

	pool.clients[creds.AccountID] = client

	return client, nil
}
