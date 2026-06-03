package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository"
)

type TokenRepo struct {
	indexTokenID map[string]*model.Token
	indexUserID  map[string][]*model.Token
	mx           sync.RWMutex
}

func NewTokenRepo() *TokenRepo {
	return &TokenRepo{
		indexTokenID: make(map[string]*model.Token),
		indexUserID:  make(map[string][]*model.Token),
	}
}

func (repo *TokenRepo) FindByID(ctx context.Context, tokenID string) (*model.Token, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	t, ok := repo.indexTokenID[tokenID]
	if ok {
		return t, nil
	}
	return nil, repository.ErrNotFound
}

func (repo *TokenRepo) ListUserTokens(ctx context.Context, userID string) ([]*model.Token, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	list, ok := repo.indexUserID[userID]
	if ok {
		return list, nil
	}
	return nil, repository.ErrNotFound
}

func (repo *TokenRepo) DeleteToken(ctx context.Context, tokenID string) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	token, err := repo.FindByID(ctx, tokenID)
	if err != nil {
		return err
	}
	delete(repo.indexUserID, token.UserID)
	delete(repo.indexTokenID, tokenID)
	return nil
}

func (repo *TokenRepo) SaveToken(ctx context.Context, token *model.Token) error {
	if token.ID == "" {
		return errors.New("token id is not defined")
	}

	repo.mx.Lock()
	defer repo.mx.Unlock()

	repo.indexTokenID[token.ID] = token
	list, ok := repo.indexUserID[token.UserID]
	if !ok {
		list = make([]*model.Token, 0, 4)
		repo.indexUserID[token.UserID] = list
	}

	list = append(list, token)
	return nil
}
