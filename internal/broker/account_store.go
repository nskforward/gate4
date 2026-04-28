package broker

import "github.com/nskforward/gate4/pkg/object"

type AccountStore struct {
	storage *object.FileStorage[*Account]
}

func NewAccountStore(filename string) (*AccountStore, error) {
	storage, err := object.NewFileStorage[*Account](filename)
	if err != nil {
		return nil, err
	}
	return &AccountStore{
		storage: storage,
	}, nil
}

func (s *AccountStore) Lookup(key string) (*Account, bool) {
	return s.storage.Get(key)
}

func (s *AccountStore) Set(key string, account *Account) error {
	return s.storage.Set(key, account)
}

func (s *AccountStore) Del(key string) error {
	return s.storage.Del(key)
}

func (s *AccountStore) Keys() []string {
	return s.storage.Keys()
}

func (s *AccountStore) Accounts() []*Account {
	return s.storage.Objects()
}
