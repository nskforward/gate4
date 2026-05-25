package finam

import (
	"context"
	"sync"
)

type Pool struct {
	ctx     context.Context
	clients map[string]*Client
	mx      sync.Mutex
}

func NewPool(ctx context.Context) *Pool {
	return &Pool{
		ctx:     ctx,
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

	client, err := newClient(pool.ctx, accountID, secret)
	if err != nil {
		return nil, err
	}

	pool.clients[accountID] = client

	return client, nil
}
