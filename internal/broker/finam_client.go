package broker

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/peers"
)

type FinamClient struct {
	clients     map[string]*finam.Client
	mx          sync.Mutex
	quoteStream *peers.PubSub[*marketdata.Quote]
}

func NewFinamClient() *FinamClient {
	var quoteStreamActive atomic.Bool

	return &FinamClient{
		clients: make(map[string]*finam.Client),
		quoteStream: peers.NewPubSub(peers.PubSubConfig[*marketdata.Quote]{
			OnStart: func(key string, group *peers.Group[*marketdata.Quote]) {
				slog.Info("start stream", "symbol", key)
				quoteStreamActive.Store(true)
				for quoteStreamActive.Load() {
					group.Send(&marketdata.Quote{
						Symbol: key,
					})
				}
			},
			OnStop: func(key string) {
				slog.Info("stop stream", "symbol", key)
				quoteStreamActive.Store(false)
			},
		}),
	}
}

func (c *FinamClient) SubscribeQuotes(account *Account, symbol string, stream pb.Admin_QuoteStreamServer) error {
	client, err := c.getClient(account)
	if err != nil {
		return err
	}
	iterator, err := client.SubscribeQuotes(stream.Context(), []string{symbol})
	if err != nil {
		return err
	}
	for q, err := range iterator {
		if err != nil {
			return err
		}
		err = stream.Send(&pb.QuoteStreamResponse{
			BrokerId:  "finam",
			Symbol:    symbol,
			Timestamp: q.Timestamp.Seconds,
			Ask: &pb.Level{
				Price: q.Ask.Value,
				Size:  q.AskSize.Value,
			},
			Bid: &pb.Level{
				Price: q.Bid.Value,
				Size:  q.BidSize.Value,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FinamClient) GetAccountInfo(ctx context.Context, account *Account) (*pb.AccountResponse, error) {
	client, err := c.getClient(account)
	if err != nil {
		return nil, err
	}
	resp, err := client.GetAccountInfo(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	return &pb.AccountResponse{
		BrokerId:  "finam",
		AccountId: resp.AccountId,
	}, nil
}

func (c *FinamClient) getClient(account *Account) (*finam.Client, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	client, ok := c.clients[account.ID]
	if !ok {
		newClient, err := finam.NewClient("api.finam.ru:443", account.Secret)
		if err != nil {
			return nil, err
		}
		c.clients[account.ID] = newClient
		client = newClient
	}

	return client, nil
}
