package finam

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/tools"
	"google.golang.org/grpc/metadata"
)

type TokenStore struct {
	ctx   context.Context
	token string
	mx    sync.RWMutex
}

func NewTokenStore(ctx context.Context) *TokenStore {
	return &TokenStore{
		ctx: ctx,
	}
}

func (store *TokenStore) RequestContext(parent context.Context) (context.Context, context.CancelFunc) {
	token, err := store.getToken()
	if err != nil {
		return parent, func() {}
	}
	mergedContext, cancel := tools.MergeContext(store.ctx, parent)
	reqCtx := metadata.AppendToOutgoingContext(mergedContext, "Authorization", token)
	return reqCtx, cancel
}

func (store *TokenStore) getToken() (string, error) {
	attempts := 0
	sleep := 10 * time.Millisecond

	for {
		attempts++
		select {
		case <-store.ctx.Done():
			return "", store.ctx.Err()

		case <-time.After(sleep):
			store.mx.RLock()
			token := store.token
			store.mx.RUnlock()

			if token != "" {
				attempts = 0
				sleep = 100 * time.Millisecond
				return token, nil
			}

			if attempts > 10 {
				return "", fmt.Errorf("too many attempts to get a new acceess token")
			}

			sleep = sleep * 2
		}
	}
}

func (store *TokenStore) set(token string) {
	store.mx.Lock()
	defer store.mx.Unlock()
	store.token = token
}
