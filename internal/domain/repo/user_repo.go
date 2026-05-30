package repo

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/domain/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type userRepo struct {
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (repo *userRepo) List(ctx context.Context) ([]model.User, error) {
	return nil, errors.New("userRepo not implemented")
}
