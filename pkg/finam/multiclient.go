package finam

import (
	"sync"
)

type MultiClient struct {
	accounts map[string]*Client
	mx       sync.Mutex
}

func NewMultiClient() *MultiClient {
	return &MultiClient{
		accounts: make(map[string]*Client),
	}
}

func (mc *MultiClient) Get(accountID, secret string) (*Client, error) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	client, ok := mc.accounts[accountID]
	if !ok {
		newClient, err := NewClient(accountID, secret)
		if err != nil {
			return nil, err
		}
		mc.accounts[accountID] = newClient
		client = newClient
	}
	return client, nil
}
