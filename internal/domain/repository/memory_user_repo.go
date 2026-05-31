package repository

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type memoryUserRepo struct {
}

func NewMemoryUserRepo() *memoryUserRepo {
	return &memoryUserRepo{}
}

func (repo *memoryUserRepo) List(ctx context.Context) ([]model.User, error) {
	return []model.User{
		{
			ID: "test_user1",
		},
	}, nil
}
