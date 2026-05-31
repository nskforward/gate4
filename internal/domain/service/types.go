package service

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
}
