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

type AccountRepo struct {
	file  string
	cache *mem.AccountRepo
	mx    sync.Mutex
}

func NewAccountRepo(cfg config.Config) (*AccountRepo, error) {
	dir, err := filepath.Abs(cfg.FileStorageDir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(dir, "accounts.json")

	repo := &AccountRepo{
		file:  file,
		cache: mem.NewAccountRepo(),
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

func (repo *AccountRepo) FindByID(ctx context.Context, accountID string) (model.Account, error) {
	return repo.cache.FindByID(ctx, accountID)
}

func (repo *AccountRepo) ListUserAccount(ctx context.Context, userID string) ([]model.Account, error) {
	return repo.cache.ListUserAccount(ctx, userID)
}

func (repo *AccountRepo) SaveAccount(ctx context.Context, account model.Account) error {
	err := repo.cache.SaveAccount(ctx, account)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *AccountRepo) DeleteAccount(ctx context.Context, accountID string) error {
	err := repo.cache.DeleteAccount(ctx, accountID)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *AccountRepo) loadFromFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Open(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Unmarshal(context.Background(), f)
}

func (repo *AccountRepo) saveToFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Create(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Marshal(context.Background(), f)
}
