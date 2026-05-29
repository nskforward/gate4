package users

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, userID string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, userID string) error
	Block(ctx context.Context, userID string, blocked bool) error
	Update(ctx context.Context, userID, secret string, expires time.Time) (*User, error)
}
