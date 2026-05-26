package finam

import (
	"context"
	"sync"
	"time"
)

type TokenStore struct {
	conn   *Conn
	token  *Token
	mx     sync.Mutex
	secret string
}

func NewTokenStore(conn *Conn, secret string) *TokenStore {
	return &TokenStore{
		conn:   conn,
		secret: secret,
	}
}

func (store *TokenStore) GetToken(ctx context.Context) (Token, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	if store.token != nil && store.token.BeforeExpiration() > 5*time.Minute && store.token.ctx.Err() == nil {
		return *store.token, nil
	}

	t, err := store.conn.Authorize(ctx, store.secret)
	if err != nil {
		return Token{}, err
	}

	store.token = &t

	return t, nil
}

func (store *TokenStore) RefreshToken(ctx context.Context, prev Token) (Token, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	if store.token != nil && store.token.created.Unix() != prev.created.Unix() {
		return *store.token, nil
	}

	t, err := store.conn.Authorize(ctx, store.secret)
	if err != nil {
		return Token{}, err
	}

	store.token = &t

	return t, nil
}
