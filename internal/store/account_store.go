package store

import (
	"fmt"
	"sync"

	"github.com/nskforward/gate4/pkg/pb"
)

type AccountStore struct {
	items map[string]*pb.Account
	mx    sync.Mutex
}

func NewAccountStore() *AccountStore {
	return &AccountStore{
		items: make(map[string]*pb.Account),
	}
}

func (s *AccountStore) List() []*pb.Account {
	s.mx.Lock()
	defer s.mx.Unlock()
	result := make([]*pb.Account, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, &pb.Account{
			Id:       item.Id,
			BrokerId: item.BrokerId,
		})
	}
	return result
}

func (s *AccountStore) Set(account *pb.Account) error {
	if account == nil || account.Id == "" || account.BrokerId == "" {
		return fmt.Errorf("input account fields must be filled")
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	s.items[account.Id] = account
	return nil
}

func (s *AccountStore) Get(id string) *pb.Account {
	s.mx.Lock()
	defer s.mx.Unlock()
	account, ok := s.items[id]
	if !ok {
		return nil
	}
	return account
}

func (s *AccountStore) Del(id string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.items, id)
}
