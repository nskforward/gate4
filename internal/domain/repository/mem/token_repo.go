package mem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository"
)

type TokenRepo struct {
	tokens map[string]model.Token
	mx     sync.RWMutex
}

func NewTokenRepo() *TokenRepo {
	return &TokenRepo{
		tokens: make(map[string]model.Token),
	}
}

func (repo *TokenRepo) FindByID(ctx context.Context, tokenID string) (model.Token, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	t, ok := repo.tokens[tokenID]
	if ok {
		return t, nil
	}
	return model.Token{}, repository.ErrNotFound
}

func (repo *TokenRepo) ListUserTokens(ctx context.Context, userID string) ([]model.Token, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	tokens := make([]model.Token, 0, 4)
	for _, t := range repo.tokens {
		if t.UserID == userID {
			tokens = append(tokens, t)
		}
	}
	return tokens, nil
}

func (repo *TokenRepo) DeleteToken(ctx context.Context, tokenID string) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	delete(repo.tokens, tokenID)
	return nil
}

func (repo *TokenRepo) SaveToken(ctx context.Context, token model.Token) error {
	if token.ID == "" {
		return errors.New("token id is not defined")
	}
	repo.mx.Lock()
	defer repo.mx.Unlock()
	repo.tokens[token.ID] = token
	return nil
}

func (repo *TokenRepo) Marshal(ctx context.Context, w io.Writer) error {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	return json.NewEncoder(w).Encode(repo.tokens)
}

func (repo *TokenRepo) Unmarshal(ctx context.Context, r io.Reader) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	return json.NewDecoder(r).Decode(&repo.tokens)
}
