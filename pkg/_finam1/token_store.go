package finam

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TokenStore struct {
	ctx         context.Context
	cancel      context.CancelFunc
	secret      string
	account     string
	token       string
	mx          sync.RWMutex
	authService auth.AuthServiceClient
	tokenCtx    context.Context
	cancelToken context.CancelFunc
}

func NewTokenStore(ctx context.Context, cancel context.CancelFunc, account, secret string, conn *grpc.ClientConn) *TokenStore {
	return &TokenStore{
		ctx:         ctx,
		cancel:      cancel,
		account:     account,
		secret:      secret,
		authService: auth.NewAuthServiceClient(conn),
	}
}

func (store *TokenStore) RequestContext(parent context.Context) (context.Context, context.CancelFunc) {
	store.mx.RLock()
	defer store.mx.RUnlock()
	ctx, cancel := tools.MergeContext(store.ctx, parent)
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "Authorization", store.token)
	return ctxWithToken, cancel
}

func (store *TokenStore) RefreshToken() error {
	store.mx.Lock()
	defer store.mx.Unlock()

	var retry retries.Retry
	for attempt := range retry.Range(store.ctx) {
		err := store.refreshToken()
		if err != nil {
			slog.Error("finam refresh token failed", "account", store.account, "attempt", attempt, "msg", err.Error())
			continue
		}
		return nil
	}
	return retry.Err()
}

func (store *TokenStore) refreshToken() error {
	ctx, cancel := context.WithTimeout(store.ctx, 30*time.Second)
	defer cancel()
	resp, err := store.authService.Auth(ctx, &auth.AuthRequest{
		Secret: store.secret,
	})
	if err != nil {
		return err
	}

	if store.cancelToken != nil {
		store.cancelToken() // close the old watch goroutine
	}

	store.token = resp.GetToken()

	resp2, err := store.authService.TokenDetails(ctx, &auth.TokenDetailsRequest{
		Token: store.token,
	})
	if err != nil {
		return err
	}

	go store.watchToken(resp2.ExpiresAt.AsTime())

	slog.Debug("refreshed finam auth token", "account", store.account, "token", TokenSuffix(resp.GetToken()))

	return nil
}

func (store *TokenStore) watchToken(expires time.Time) {
	store.tokenCtx, store.cancelToken = context.WithCancel(store.ctx)

	slog.Debug("wait finam token expiration", "account", store.account, "expiration", expires)

	wait := time.Until(expires) - 10*time.Minute

	select {
	case <-store.tokenCtx.Done():
		slog.Debug("finam watch token cancelled", "account", store.account)
		return
	case <-time.After(wait):
		err := store.RefreshToken()
		if err != nil {
			slog.Error("cannot refresh finam token", "account", store.account, "reason", err.Error())
			store.cancel()
		}
	}
}

func TokenSuffix(token string) string {
	if len(token) > 8 {
		return token[len(token)-8:]
	}
	return token
}
