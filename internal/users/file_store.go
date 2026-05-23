package users

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/tools"
)

type FileStore struct {
	filename string
	users    map[string]*User
	mx       sync.RWMutex
}

func NewFileStorage() *FileStore {
	normalized, err := tools.NormalizeFilename("data/users.json")
	if err != nil {
		panic(fmt.Errorf("bad users file storage filename path: %w", err))
	}

	store := &FileStore{
		filename: normalized,
		users:    make(map[string]*User),
	}

	err = store.load()
	if err != nil {
		panic(fmt.Errorf("cannot load users from file storage: %w", err))
	}
	return store
}

func (store *FileStore) List(ctx context.Context) ([]*User, error) {
	store.mx.RLock()
	defer store.mx.RUnlock()
	users := make([]*User, 0, len(store.users))
	for _, user := range store.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Created.Before(users[i].Created)
	})
	return users, nil
}

func (store *FileStore) Find(ctx context.Context, userID string) (*User, error) {
	store.mx.RLock()
	defer store.mx.RUnlock()
	user, ok := store.users[userID]
	if ok {
		return user, nil
	}
	return nil, nil
}

func (store *FileStore) Create(ctx context.Context, user *User) error {
	store.mx.Lock()
	defer store.mx.Unlock()
	user.ID = store.generateID(user)
	user.Created = time.Now()
	store.users[user.ID] = user
	return store.save()
}

func (store *FileStore) Delete(ctx context.Context, userID string) error {
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.users, userID)
	return store.save()
}

func (store *FileStore) Block(ctx context.Context, userID string, blocked bool) error {
	user, _ := store.Find(ctx, userID)
	if user == nil {
		return fmt.Errorf("user not found with id: %s", userID)
	}
	store.mx.Lock()
	defer store.mx.Unlock()
	user.Blocked = blocked
	return store.save()
}

func (store *FileStore) Update(ctx context.Context, userID, secret string, expires time.Time) (*User, error) {
	user, _ := store.Find(ctx, userID)
	if user == nil {
		return nil, fmt.Errorf("user not found with id: %s", userID)
	}
	store.mx.Lock()
	defer store.mx.Unlock()
	user.Secret = secret
	user.Expires = expires
	err := store.save()
	return user, err
}

func (store *FileStore) generateID(user *User) string {
	return fmt.Sprintf("%s.%s", string(user.BrokerID), user.AccountID)
}

func (store *FileStore) load() error {
	f, err := os.Open(store.filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&store.users)
}

func (store *FileStore) save() error {
	f, err := os.OpenFile(store.filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(store.users)
}
