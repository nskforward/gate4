package account

import "github.com/nskforward/gate4/pkg/object"

type Store struct {
	storage *object.FileStorage[*Account]
}

func NewStore(filename string) (*Store, error) {
	storage, err := object.NewFileStorage[*Account](filename)
	if err != nil {
		return nil, err
	}
	return &Store{
		storage: storage,
	}, nil
}

func (s *Store) Lookup(key string) (*Account, bool) {
	return s.storage.Get(key)
}

func (s *Store) Set(key string, account *Account) error {
	return s.storage.Set(key, account)
}

func (s *Store) Del(key string) error {
	return s.storage.Del(key)
}

func (s *Store) Keys() []string {
	return s.storage.Keys()
}

func (s *Store) Accounts() []*Account {
	return s.storage.Objects()
}
