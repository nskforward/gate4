package finam

import (
	"context"
	"sync"

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
	ctx, cancel := tools.MergeContext(store.ctx, parent)
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "Authorization", store.Get())
	return ctxWithToken, cancel
}

func (store *TokenStore) Get() string {
	store.mx.RLock()
	defer store.mx.RUnlock()
	return store.token
}

func (store *TokenStore) Set(token string) {
	store.mx.Lock()
	defer store.mx.Unlock()
	store.token = token
}

func TokenSuffix(token string) string {
	if len(token) > 8 {
		return token[len(token)-8:]
	}
	return token
}
