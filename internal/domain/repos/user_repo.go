package repos

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type userRepo struct {
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (repo *userRepo) List(ctx context.Context) ([]domain.User, error) {
	return nil, errors.New("userRepo not implemented")
}
