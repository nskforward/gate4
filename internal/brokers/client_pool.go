package brokers

import (
	"errors"

	"github.com/nskforward/gate4/internal/domain/users"
	"github.com/nskforward/gate4/pkg/finam"
)

type ClientPool struct {
	finamClients *finam.ClientPool
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		finamClients: finam.NewClientPool(),
	}
}

func (pool *ClientPool) GetOrCreate(user *users.User) (Client, error) {
	if user.BrokerID == users.FINAM {
		return pool.finamClients.GetOrCreateClient(&finam.Creds{AccountID: user.AccountID, Secret: user.Secret})
	}
	return nil, errors.New("unknown user broker")
}

func (pool *ClientPool) Delete(user *users.User) error {
	if user.BrokerID == users.FINAM {
		return pool.finamClients.DeleteClient(user.AccountID)
	}
	return errors.New("unknown user broker")
}
