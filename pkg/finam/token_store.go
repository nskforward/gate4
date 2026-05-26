package finam

import (
	"context"
	"fmt"
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

func (store *TokenStore) GetToken(ctx context.Context, minLifetime time.Duration) (Token, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	if store.token != nil && store.token.BeforeExpiration() > minLifetime && store.token.ctx.Err() == nil {
		return *store.token, nil
	}

	t, err := store.conn.Authorize(ctx, store.secret)
	if err != nil {
		return Token{}, err
	}

	store.token = &t

	fmt.Println("token expiration:", t.BeforeExpiration())

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

	fmt.Println("token expiration:", t.BeforeExpiration())

	return t, nil
}
