package repository

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/domain/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Save(ctx context.Context, newUser model.User) error
	Delete(ctx context.Context, userID string) error
}
