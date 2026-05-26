package finam

import (
	"context"
	"log/slog"
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

func (pool *Pool) Get(account, secret string) (*Client, error) {
	pool.mx.Lock()
	defer pool.mx.Unlock()

	client, ok := pool.clients[account]
	if ok {
		return client, nil
	}

	client, err := NewClient(pool.ctx, account, secret)
	if err != nil {
		return nil, err
	}

	pool.clients[account] = client

	slog.Debug("finam client created", "account", account)

	return client, nil
}
