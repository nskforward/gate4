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

type TokenRepo struct {
	file  string
	cache *mem.TokenRepo
	mx    sync.Mutex
}

func NewTokenRepo(cfg config.Config) (*TokenRepo, error) {
	dir, err := filepath.Abs(cfg.FileStorageDir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(dir, "tokens.json")

	repo := &TokenRepo{
		file:  file,
		cache: mem.NewTokenRepo(),
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

func (repo *TokenRepo) FindByID(ctx context.Context, tokenID string) (model.Token, error) {
	return repo.cache.FindByID(ctx, tokenID)
}

func (repo *TokenRepo) ListUserTokens(ctx context.Context, userID string) ([]model.Token, error) {
	return repo.cache.ListUserTokens(ctx, userID)
}

func (repo *TokenRepo) SaveToken(ctx context.Context, token model.Token) error {
	err := repo.cache.SaveToken(ctx, token)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *TokenRepo) DeleteToken(ctx context.Context, tokenID string) error {
	err := repo.cache.DeleteToken(ctx, tokenID)
	if err != nil {
		return err
	}
	return repo.saveToFile()
}

func (repo *TokenRepo) loadFromFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Open(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Unmarshal(context.Background(), f)
}

func (repo *TokenRepo) saveToFile() error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	f, err := os.Create(repo.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return repo.cache.Marshal(context.Background(), f)
}
