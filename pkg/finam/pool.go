package finam

import (
	"sync"
)

type Pool struct {
	clients map[string]*Client
	mx      sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		clients: make(map[string]*Client),
	}
}

func (pool *Pool) Get(accountID, secret string) (*Client, error) {
	pool.mx.Lock()
	defer pool.mx.Unlock()

	client, ok := pool.clients[accountID]
	if ok {
		return client, nil
	}

	client, err := newClient(accountID, secret)
	if err != nil {
		return nil, err
	}

	pool.clients[accountID] = client

	return client, nil
}
