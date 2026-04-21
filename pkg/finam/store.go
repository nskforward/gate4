package finam

import (
	"fmt"
	"sync"
)

type Store struct {
	addr     string
	accounts map[string]*Leaf
	mx       sync.RWMutex
}

type Leaf struct {
	secret string
	client *Client
	mx     sync.Mutex
}

func NewStore(addr string) *Store {
	return &Store{
		addr:     addr,
		accounts: make(map[string]*Leaf),
	}
}

func (s *Store) Set(accountID, secret string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.accounts[accountID] = &Leaf{
		secret: secret,
	}
}

func (s *Store) Get(accountID string) (*Client, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	leaf, ok := s.accounts[accountID]
	if !ok {
		return nil, fmt.Errorf("account %s has no stored secret", accountID)
	}

	leaf.mx.Lock()
	defer leaf.mx.Unlock()

	if leaf.client != nil {
		return leaf.client, nil
	}

	client, err := NewClient(s.addr, leaf.secret)
	if err != nil {
		return nil, err
	}

	leaf.client = client

	return client, nil
}
