package broker

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type Broker struct {
	accountStore  *AccountStore
	positionStore *PositionStore
	finamClient   *FinamClient
}

func NewBroker() (*Broker, error) {
	storage, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	b := &Broker{
		accountStore:  storage,
		positionStore: NewPositionStore(),
		finamClient:   NewFinamClient(),
	}
	err = b.importPositions()
	if err != nil {
		return nil, err
	}
	return b, nil
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

func (b *Broker) GetPositions(account *Account) []*pb.Position {
	return b.positionStore.Get(account)
}

func (b *Broker) SubscribeQuotes(req *pb.QuoteStreamRequest, stream grpc.ServerStreamingServer[pb.QuoteStreamResponse]) error {
	account, err := b.lookupAccount(req.AccountKey)
	if err != nil {
		return err
	}
	client, err := b.LookupClient(account)
	if err != nil {
		return err
	}
	return client.SubscribeQuotes(account, req.Symbol, stream)
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

func (b *Broker) LookupAccount(key string) *Account {
	account, ok := b.accountStore.Lookup(key)
	if !ok {
		return nil
	}
	return account
}

func (b *Broker) LookupClient(account *Account) (Client, error) {
	switch account.Broker {
	case "finam":
		return b.finamClient, nil

	default:
		return nil, fmt.Errorf("unknown broker: %s", account.Broker)
	}
}

func (b *Broker) UpdatePositions(account *Account, positions []*pb.Position) {
	b.positionStore.Update(account, positions)
}

func (b *Broker) GetSchedule(ctx context.Context, account *Account, symbol string) ([]*pb.ScheduleSession, *pb.ScheduleSession, error) {
	sessions, current, err := b.finamClient.Schedule(ctx, account, symbol)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*pb.ScheduleSession, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, &pb.ScheduleSession{
			Type:  string(sess.Type),
			Start: sess.Start,
			End:   sess.End,
		})
	}
	return result, &pb.ScheduleSession{
		Type:  string(current.Type),
		Start: current.Start,
		End:   current.End,
	}, nil
}

func (b *Broker) lookupAccount(key string) (*Account, error) {
	account, ok := b.accountStore.Lookup(key)
	if !ok {
		return account, fmt.Errorf("unknown account key: %s", key)
	}
	return account, nil
}

func (b *Broker) importPositions() error {
	accounts := b.accountStore.Accounts()
	for _, account := range accounts {
		client, err := b.LookupClient(account)
		if err != nil {
			return err
		}
		resp, err := client.GetAccountInfo(context.Background(), account)
		if err != nil {
			return err
		}
		b.UpdatePositions(account, resp.Positions)
	}
	return nil
}
