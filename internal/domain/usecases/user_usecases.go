package usecases

import (
	"context"

	"github.com/nskforward/gate4/internal/domain"
)

type UserRepo interface {
	List(ctx context.Context) ([]domain.User, error)
}

type UserUsecases struct {
	userRepo UserRepo
}

func NewUserUsecases(userRepo UserRepo) *UserUsecases {
	return &UserUsecases{
		userRepo: userRepo,
	}
}

func (usecase *UserUsecases) List(ctx context.Context) ([]domain.User, error) {
	return usecase.userRepo.List(ctx)
}
