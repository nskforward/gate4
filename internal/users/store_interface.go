package users

import (
	"context"
	"time"
)

type Store interface {
	Add(ctx context.Context, user *User) error
	Find(ctx context.Context, userID string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, userID string) error
	Block(ctx context.Context, userID string, blocked bool) error
	Update(ctx context.Context, userID, secret string, validUntil time.Time) error
}
