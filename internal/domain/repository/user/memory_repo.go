package user

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/domain/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type memoryRepo struct {
}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{}
}

func (repo *memoryRepo) List(ctx context.Context) ([]model.User, error) {
	return []model.User{
		{
			ID: "test_user1",
		},
	}, nil
}
