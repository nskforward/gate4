package broker

import (
	"fmt"
	"sync"

	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/types"
)

type Broker struct {
	accounts     *AccountStore
	finamClients *finam.MultiClient
	mx           sync.Mutex
}

func NewBroker() (*Broker, error) {
	accounts, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	return &Broker{
		accounts:     accounts,
		finamClients: finam.NewMultiClient(),
	}, nil
}

func (b *Broker) AddAccount(account *Account) error {
	return b.accounts.Set(account.Key(), account)
}

func (b *Broker) DelAccount(key string) error {
	return b.accounts.Del(key)
}

func (b *Broker) Accounts() []*Account {
	return b.accounts.Accounts()
}

func (b *Broker) Client(key string) (types.BrokerClient, error) {
	account, ok := b.accounts.Lookup(key)
	if !ok {
		return nil, fmt.Errorf("unknown account key: %s", key)
	}
	switch account.Broker {
	case "finam":
		return b.finamClients.Get(account.ID, account.Secret)
	}
	return nil, fmt.Errorf("client not registered for broker: %s", account.Broker)
}
