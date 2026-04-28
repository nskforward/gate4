package broker

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/pkg/pb"
)

type Broker struct {
	accountStore *AccountStore
	finamClient  *FinamClient
}

func NewBroker() (*Broker, error) {
	storage, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	return &Broker{
		accountStore: storage,
		finamClient:  NewFinamClient(),
	}, nil
}

func (b *Broker) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	account, err := b.lookupAccount(in.AccountKey)
	if err != nil {
		return nil, err
	}
	client, err := b.lookupClient(account)
	if err != nil {
		return nil, err
	}
	return client.GetAccountInfo(ctx, account)
}

func (b *Broker) lookupAccount(key string) (Account, error) {
	account, ok := b.accountStore.Lookup(key)
	if !ok {
		return account, fmt.Errorf("unknown account key: %s", key)
	}
	return account, nil
}

func (b *Broker) lookupClient(account Account) (Client, error) {
	switch account.Broker {
	case "finam":
		return b.finamClient, nil

	default:
		return nil, fmt.Errorf("unknown broker: %s", account.Broker)
	}
}
