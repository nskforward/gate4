package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/repository"
)

type UserService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Update updates all fields except ID and Created.
func (s *UserService) Update(ctx context.Context, newUser model.User) error {
	err := newUser.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	oldUser, err := s.userRepo.FindByID(ctx, newUser.ID)
	if err != nil {
		return err
	}
	oldUser.Email = newUser.Email
	oldUser.Name = newUser.Name
	oldUser.Blocked = newUser.Blocked
	oldUser.Role = newUser.Role
	return s.userRepo.Save(ctx, oldUser)
}

func (s *UserService) FindByID(ctx context.Context, userID string) (model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	return s.userRepo.List(ctx)
}

func (s *UserService) Create(ctx context.Context, newUser *model.User) error {
	err := newUser.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	newUser.ID = uuid.NewString()
	newUser.Created = time.Now()
	return s.userRepo.Save(ctx, *newUser)
}

func (s *UserService) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id must be specified")
	}

	tokens, err := s.tokenRepo.ListUserTokens(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		return errors.New("cannot delete user with active tokens")
	}

	return s.userRepo.Delete(ctx, userID)
}
