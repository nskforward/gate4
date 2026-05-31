package repository

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/domain/model"
)

type userMemoryRepo struct {
}

func NewUserMemoryRepo() *userMemoryRepo {
	return &userMemoryRepo{}
}

func (repo *userMemoryRepo) List(ctx context.Context) ([]model.User, error) {
	return []model.User{
		{
			ID:      "test_user1",
			Created: time.Now(),
		},
	}, nil
}
