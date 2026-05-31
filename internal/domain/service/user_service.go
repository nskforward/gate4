package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	return s.userRepo.List(ctx)
}

func (s *UserService) Create(ctx context.Context, newUser *model.User) error {
	if newUser.ID != "" {
		return errors.New("user id must be empty for a new user")
	}
	if newUser.Name == "" {
		return errors.New("user name cannot be empty")
	}
	if newUser.Email == "" {
		return errors.New("user email cannot be empty")
	}
	newUser.ID = uuid.NewString()
	newUser.Created = time.Now()
	return s.userRepo.Create(ctx, *newUser)
}

func (s *UserService) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id must be specified")
	}
	return s.userRepo.Delete(ctx, userID)
}
