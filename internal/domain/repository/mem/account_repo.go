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

type AccountRepo struct {
	accounts map[string]model.Account
	mx       sync.RWMutex
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{
		accounts: make(map[string]model.Account),
	}
}

func (repo *AccountRepo) FindByID(ctx context.Context, accountID string) (model.Account, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	t, ok := repo.accounts[accountID]
	if ok {
		return t, nil
	}
	return model.Account{}, repository.ErrNotFound
}

func (repo *AccountRepo) ListUserAccount(ctx context.Context, userID string) ([]model.Account, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	accounts := make([]model.Account, 0, 4)
	for _, t := range repo.accounts {
		if t.UserID == userID {
			accounts = append(accounts, t)
		}
	}
	return accounts, nil
}

func (repo *AccountRepo) SaveAccount(ctx context.Context, account model.Account) error {
	if account.ID == "" {
		return errors.New("account id is not defined")
	}
	repo.mx.Lock()
	defer repo.mx.Unlock()
	repo.accounts[account.ID] = account
	return nil
}

func (repo *AccountRepo) DeleteAccount(ctx context.Context, accountID string) error {
	account, err := repo.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	repo.mx.Lock()
	defer repo.mx.Unlock()
	delete(repo.accounts, account.ID)
	return nil
}

func (repo *AccountRepo) Marshal(ctx context.Context, w io.Writer) error {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	return json.NewEncoder(w).Encode(repo.accounts)
}

func (repo *AccountRepo) Unmarshal(ctx context.Context, r io.Reader) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	return json.NewDecoder(r).Decode(&repo.accounts)
}
