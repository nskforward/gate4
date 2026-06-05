package mem

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"slices"
	"sync"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository"
)

type UserRepo struct {
	users map[string]model.User
	mx    sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]model.User, 8),
	}
}

func (repo *UserRepo) FindByID(ctx context.Context, userID string) (model.User, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	user, ok := repo.users[userID]
	if !ok {
		return model.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (repo *UserRepo) List(ctx context.Context) ([]model.User, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	return slices.Collect(maps.Values(repo.users)), nil
}

func (repo *UserRepo) Save(ctx context.Context, newUser model.User) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	repo.users[newUser.ID] = newUser
	return nil
}

func (repo *UserRepo) Delete(ctx context.Context, userID string) error {
	user, err := repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	repo.mx.Lock()
	defer repo.mx.Unlock()
	delete(repo.users, user.ID)
	return nil
}

func (repo *UserRepo) Marshal(ctx context.Context, w io.Writer) error {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	return json.NewEncoder(w).Encode(repo.users)
}

func (repo *UserRepo) Unmarshal(ctx context.Context, r io.Reader) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	return json.NewDecoder(r).Decode(&repo.users)
}
