package usecase

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type UserRepo interface {
	List(ctx context.Context) ([]model.User, error)
}

type UserUsecase struct {
	userRepo UserRepo
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (usecase *UserUsecase) List(ctx context.Context) ([]model.User, error) {
	return usecase.userRepo.List(ctx)
}
