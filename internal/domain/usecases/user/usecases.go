package user

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
}

type UserUsecases struct {
	userRepo UserRepository
}

func NewUserUsecases(userRepo UserRepository) *UserUsecases {
	return &UserUsecases{
		userRepo: userRepo,
	}
}

func (u *UserUsecases) List(ctx context.Context) ([]model.User, error) {
	return u.userRepo.List(ctx)
}
