package service

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, newUser model.User) error
	Delete(ctx context.Context, userID string) error
}
