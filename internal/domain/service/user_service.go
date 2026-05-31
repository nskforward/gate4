package service

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) List(ctx context.Context) ([]model.User, error) {
	return u.userRepo.List(ctx)
}
