package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrKeyEmpty = errors.New("key is empty")
)

type ObjectStorage interface {

	// Keys returns list of all keys
	Keys(ctx context.Context) []string

	// Save creates or replaces existing object
	// throws: ErrKeyEmpty
	Save(ctx context.Context, key string, src any) error

	// Find finds an existing object
	// throws: ErrKeyEmpty | ErrNotFound
	Find(ctx context.Context, key string, dst any) error

	// Delete deletes only the existing data
	// throws: ErrKeyEmpty | ErrNotFound
	Delete(ctx context.Context, key string) error
}
