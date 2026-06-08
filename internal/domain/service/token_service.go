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

type TokenService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewTokenService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *TokenService {
	return &TokenService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *TokenService) CreateToken(ctx context.Context, token *model.Token) error {
	token.ID = uuid.NewString()
	token.Created = time.Now()
	err := token.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return fmt.Errorf("cannot find user by id: %w", err)
	}

	if user.Blocked {
		return errors.New("cannot create a token for blocked user")
	}

	return s.tokenRepo.SaveToken(ctx, *token)
}

func (s *TokenService) ListUserTokens(ctx context.Context, userID string) ([]model.Token, error) {
	return s.tokenRepo.ListUserTokens(ctx, userID)
}

func (s *TokenService) DeleteToken(ctx context.Context, tokenID string) error {
	return s.tokenRepo.DeleteToken(ctx, tokenID)
}

func (s *TokenService) FindUser(ctx context.Context, tokenID string) (model.User, error) {

	token, err := s.tokenRepo.FindByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.User{}, fmt.Errorf("unknown token")
		}
		return model.User{}, err
	}

	err = token.Validate()
	if err != nil {
		return model.User{}, fmt.Errorf("token validation failed: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.User{}, fmt.Errorf("unknown token user")
		}
		return model.User{}, err
	}

	return user, nil
}
