package store

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
)

type AccountStore struct {
	items    map[string]*pb.Account
	mx       sync.RWMutex
	provider AccountProvider
}

func NewAccountStore(provider AccountProvider) (*AccountStore, error) {
	items, err := provider.Load()
	if err != nil {
		return nil, err
	}
	s := &AccountStore{
		items:    items,
		provider: provider,
	}
	return s, nil
}

func (s *AccountStore) List() []*pb.Account {
	s.mx.RLock()
	defer s.mx.RUnlock()
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
	err := validateAccount(account)
	if err != nil {
		return fmt.Errorf("account validation failed: %w", err)
	}
	key := fmt.Sprintf("%s.%s", account.BrokerId, account.Id)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.items[key] = account
	return s.provider.Save(s.items)
}

func (s *AccountStore) Get(id string) *pb.Account {
	s.mx.RLock()
	defer s.mx.RUnlock()
	account, ok := s.items[id]
	if !ok {
		return nil
	}
	return account
}

func (s *AccountStore) Del(id string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.items, id)
	return s.provider.Save(s.items)
}

func validateAccount(account *pb.Account) error {
	if account == nil {
		return fmt.Errorf("account cannot be a nil")
	}
	if account.Id == "" {
		return fmt.Errorf("account id cannot be empty")
	}
	if !slices.Contains([]string{"finam"}, account.BrokerId) {
		return fmt.Errorf("unknown broker id: %s", account.BrokerId)
	}
	if account.ValidUntil < time.Now().Unix() {
		return fmt.Errorf("valid date year cannot be in the past")
	}
	return nil
}
