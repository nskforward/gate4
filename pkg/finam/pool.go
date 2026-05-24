package finam

import "errors"

type Pool struct {
}

func NewPool() *Pool {
	return &Pool{}
}

func (pool *Pool) Get(accountID, secret string) (*Client, error) {
	return nil, errors.New("finam pool: not implemented")
}
