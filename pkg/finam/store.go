package finam

import (
	"sync"
)

type Store struct {
	addr     string
	accounts map[string]*Client
	mx       sync.Mutex
}

func NewStore(addr string) *Store {
	return &Store{
		addr:     addr,
		accounts: make(map[string]*Client),
	}
}

func (s *Store) GetClient(accountID, secret string) (*Client, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	client, ok := s.accounts[accountID]
	if ok {
		err := client.UpdateSecret(secret)
		return client, err
	}

	client, err := NewClient(s.addr, secret)
	if err != nil {
		return nil, err
	}
	s.accounts[accountID] = client
	return client, nil
}
