package repository

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/domain/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Save(ctx context.Context, newUser model.User) error
	Delete(ctx context.Context, userID string) error
}

type TokenRepository interface {
	FindByID(ctx context.Context, tokenID string) (model.Token, error)
	ListUserTokens(ctx context.Context, userID string) ([]model.Token, error)
	SaveToken(ctx context.Context, token model.Token) error
	DeleteToken(ctx context.Context, tokenID string) error
}

type AccountRepository interface {
	FindByID(ctx context.Context, accountID string) (model.Account, error)
	ListUserAccount(ctx context.Context, userID string) ([]model.Account, error)
	SaveAccount(ctx context.Context, account model.Account) error
	DeleteAccount(ctx context.Context, accountID string) error
}
