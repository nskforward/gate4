package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository"
)

type AccountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, account *model.Account) error {
	account.ID = uuid.NewString()
	account.Created = time.Now()

	err := account.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	_, err = s.accountRepo.FindByID(ctx, account.UserID)
	if err != nil {
		return fmt.Errorf("cannot find user by id: %w", err)
	}

	return s.accountRepo.SaveAccount(ctx, *account)
}

func (s *AccountService) ListUserAccounts(ctx context.Context, userID string) ([]model.Account, error) {
	return s.accountRepo.ListUserAccount(ctx, userID)
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountID string) error {
	return s.accountRepo.DeleteAccount(ctx, accountID)
}
