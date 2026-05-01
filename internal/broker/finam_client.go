package broker

import (
	"context"
	"iter"
	"log/slog"
	"sync"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/peers"
	"github.com/nskforward/gate4/pkg/retries"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type FinamClient struct {
	addr        string
	clients     map[string]*finam.Client
	mx          sync.Mutex
	quoteStream *peers.PubSub[*marketdata.Quote]
}

func NewFinamClient(addr string) *FinamClient {
	return &FinamClient{
		addr:        addr,
		clients:     make(map[string]*finam.Client),
		quoteStream: peers.NewPubSub[*marketdata.Quote](),
	}
}

func (c *FinamClient) SubscribeQuotes(account *Account, symbol string, stream pb.Admin_QuoteStreamServer) error {
	slog.Debug("client connected for quote stream", "symbol", symbol)

	client, err := c.getClient(account)
	if err != nil {
		return err
	}
	group, loaded := c.quoteStream.LoadOrCreate(symbol)
	if !loaded {
		go c.subscribeQuotes(symbol, client, group)
	}

	p := group.NewPeer()
	defer p.Close()

	for {
		q, ok := p.Read(stream.Context())
		if !ok {
			slog.Debug("client disconnected from quote stream", "symbol", symbol)
			break // client disconnected
		}

		err := stream.Send(&pb.QuoteStreamResponse{
			BrokerId:  "finam",
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp.Seconds,
			Ask:       q.Ask.Value,
			Bid:       q.Bid.Value,
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

func (c *FinamClient) subscribeQuotes(symbol string, client *finam.Client, group *peers.Group[*marketdata.Quote]) {

	slog.Debug("start external quote stream", "symbol", symbol)

	retry := retries.NewRetry(
		retries.DefaultConfig(),
		func() (iter.Seq2[*marketdata.Quote, error], error) {
			return client.SubscribeQuotes(
				context.Background(),
				[]string{symbol},
			)
		},
	)

MAIN_LOOP:
	for {
		stream, err := retry.Do(context.Background())
		if err != nil {
			slog.Error("cannot subscribe for quote stream", "symbol", symbol, "error", err.Error())
			break MAIN_LOOP
		}

		ask := &decimal.Decimal{Value: "0"}
		bid := &decimal.Decimal{Value: "0"}

		for q, err := range stream {
			if err != nil {
				slog.Warn("quote stream disconnected", "symbol", symbol, "error", err.Error())
				break
			}
			modified := false
			if q.Ask != nil && q.Ask.Value != ask.Value {
				modified = true
				ask.Value = q.Ask.Value
			}
			if q.Bid != nil && q.Bid.Value != bid.Value {
				modified = true
				bid.Value = q.Bid.Value
			}
			if !modified {
				continue
			}
			q.Ask = ask
			q.Bid = bid
			if !group.Send(q) {
				slog.Warn("no subscribers for quote stream", "symbol", symbol)
				break MAIN_LOOP
			}
		}
	}

	slog.Info("exit external quote stream")
	group.Close()
}

func (c *FinamClient) getClient(account *Account) (*finam.Client, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	client, ok := c.clients[account.ID]
	if !ok {
		newClient, err := finam.NewClient(c.addr, account.Secret)
		if err != nil {
			return nil, err
		}
		c.clients[account.ID] = newClient
		client = newClient
	}

	return client, nil
}
