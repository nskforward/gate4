package brokers

import (
	"errors"

	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/finam"
)

type Pool struct {
	finamClients *finam.Pool
}

func NewPool() *Pool {
	return &Pool{
		finamClients: finam.NewPool(),
	}
}

func (pool *Pool) Get(user *users.User) (Client, error) {
	if user.BrokerID == users.FINAM {
		return pool.finamClients.Get(user.AccountID, user.Secret)
	}
	return nil, errors.New("unknown user broker")
}
