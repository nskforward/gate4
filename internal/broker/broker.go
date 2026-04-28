package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/store"
	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
)

type Broker struct {
	accountStore *store.AccountStore
	finamStore   *finam.Store
}

func NewBroker(accountStore *store.AccountStore, finamStore *finam.Store) *Broker {
	return &Broker{
		accountStore: accountStore,
		finamStore:   finamStore,
	}
}

func (b *Broker) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	account := b.accountStore.Get(in.AccountKey)
	if account == nil {
		return nil, errors.New("unknown account key")
	}

	switch account.BrokerId {
	case "finam":
		client, err := b.finamStore.GetClient(account.Id, account.Secret)
		if err != nil {
			return nil, fmt.Errorf("cannot get finam client: %w", err)
		}
		resp, err := client.GetAccountInfo(ctx, account.Id)
		if err != nil {
			return nil, fmt.Errorf("finam communication error: %w", err)
		}
		return &pb.AccountResponse{
			BrokerId:  "finam",
			AccountId: resp.AccountId,
		}, nil
	}

	return nil, fmt.Errorf("broker '%s' not supported", account.BrokerId)
}
