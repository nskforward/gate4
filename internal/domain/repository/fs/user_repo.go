package fs

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository/mem"
)

type UserRepo struct {
	file  string
	cache *mem.UserRepo
	mx    sync.Mutex
}

func NewUserRepo(cfg config.Config) (*UserRepo, error) {
	dir, err := filepath.Abs(cfg.FileStorageDir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(dir, "users.json")

	repo := &UserRepo{
		file:  file,
		cache: mem.NewUserRepo(),
	}

	_, err = os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return repo, nil
		}
		return nil, err
	}

	return repo, repo.loadFromFile()
}

func (repo *UserRepo) FindByID(ctx context.Context, userID string) (model.User, error) {
	return repo.cache.FindByID(ctx, userID)
}

func (repo *UserRepo) List(ctx context.Context) ([]model.User, error) {
	return repo.cache.List(ctx)
}

func (repo *UserRepo) Save(ctx context.Context, newUser model.User) error {
	err := repo.cache.Save(ctx, newUser)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *UserRepo) Delete(ctx context.Context, userID string) error {
	err := repo.cache.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *UserRepo) loadFromFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Open(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Unmarshal(context.Background(), f)
}

func (repo *UserRepo) saveToFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Create(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Marshal(context.Background(), f)
}
