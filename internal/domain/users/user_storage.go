package users

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/storage"
)

type UserStorage struct {
	objectStorage storage.ObjectStorage
	users         map[string]User
	mx            sync.Mutex
}

func NewUserStorage(objectStorage storage.ObjectStorage) *UserStorage {
	s := &UserStorage{
		objectStorage: objectStorage,
		users:         make(map[string]User),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.fill(ctx)
	if err != nil {
		slog.Error("cannot fill in the user storage", "reason", err.Error())
		os.Exit(1)
	}

	return s
}

func (s *UserStorage) List(ctx context.Context) ([]User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *UserStorage) Get(ctx context.Context, id string) (User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	user, ok := s.users[id]
	if !ok {
		return User{}, storage.ErrNotFound
	}
	return user, nil
}

func (s *UserStorage) Create(ctx context.Context, user *User) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	user.ID = s.generateID(user)
	user.Created = time.Now()
	s.users[user.ID] = *user
	return s.objectStorage.Save(ctx, user.ID, user)
}

func (s *UserStorage) Delete(ctx context.Context, id string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	err := s.objectStorage.Delete(ctx, id)
	if err != nil {
		return err
	}
	delete(s.users, id)
	return nil
}

func (s *UserStorage) Update(ctx context.Context, user User) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	err := s.objectStorage.Save(ctx, user.ID, user)
	if err != nil {
		return err
	}

	s.users[user.ID] = user
	return nil
}

func (s *UserStorage) generateID(user *User) string {
	return fmt.Sprintf("%s.%s", string(user.BrokerID), user.AccountID)
}

func (s *UserStorage) fill(ctx context.Context) error {
	s.users = make(map[string]User)
	keys := s.objectStorage.Keys(ctx)
	for _, key := range keys {
		var user User
		err := s.objectStorage.Read(ctx, key, &user)
		if err != nil {
			return err
		}
		s.users[user.ID] = user
	}
	return nil
}
