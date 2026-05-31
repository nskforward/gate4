package repository

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/nskforward/gate4/internal/domain/model"
)

type inMemUserRepo struct {
	users map[string]model.User
	mx    sync.RWMutex
}

func NewInMemUserRepo() *inMemUserRepo {
	return &inMemUserRepo{
		users: make(map[string]model.User, 8),
	}
}

func (repo *inMemUserRepo) List(ctx context.Context) ([]model.User, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	return slices.Collect(maps.Values(repo.users)), nil
}

func (repo *inMemUserRepo) Create(ctx context.Context, newUser model.User) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	repo.users[newUser.ID] = newUser
	return nil
}

func (repo *inMemUserRepo) Delete(ctx context.Context, userID string) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	delete(repo.users, userID)
	return nil
}
