package brokers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/finam"
)

type Pool struct {
	finamClients *finam.Pool
}

func NewPool(ctx context.Context) *Pool {
	return &Pool{
		finamClients: finam.NewPool(ctx),
	}
}

func (pool *Pool) Get(user *users.User) (Client, error) {
	if user.BrokerID == users.FINAM {
		client, loaded, err := pool.finamClients.Get(user.AccountID, user.Secret)
		if err != nil {
			return nil, err
		}
		if !loaded {
			slog.Debug("created a new finam client", "account", user.AccountID)
		}
		return client, nil
	}
	return nil, errors.New("unknown user broker")
}
