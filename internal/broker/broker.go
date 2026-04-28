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

func (b *Broker) AccountKeys() []string {
	return b.accountStore.Keys()
}

func (b *Broker) Accounts() []*Account {
	return b.accountStore.Accounts()
}

func (b *Broker) AddAccount(account *Account) error {
	return b.accountStore.Set(account.Key(), account)
}

func (b *Broker) DeleteAccount(key string) error {
	return b.accountStore.Del(key)
}

func (b *Broker) GetAccountInfo(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	account, err := b.lookupAccount(in.AccountKey)
	if err != nil {
		return nil, err
	}
	client, err := b.LookupClient(account)
	if err != nil {
		return nil, err
	}
	return client.GetAccountInfo(ctx, account)
}

func (b *Broker) lookupAccount(key string) (*Account, error) {
	account, ok := b.accountStore.Lookup(key)
	if !ok {
		return account, fmt.Errorf("unknown account key: %s", key)
	}
	return account, nil
}

func (b *Broker) LookupClient(account *Account) (Client, error) {
	switch account.Broker {
	case "finam":
		return b.finamClient, nil

	default:
		return nil, fmt.Errorf("unknown broker: %s", account.Broker)
	}
}
