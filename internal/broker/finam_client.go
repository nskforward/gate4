package broker

import (
	"context"
	"log/slog"
	"sync"
	"time"

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
	return &FinamClient{
		clients:     make(map[string]*finam.Client),
		quoteStream: peers.NewPubSub[*marketdata.Quote](),
	}
}

func (c *FinamClient) SubscribeQuotes(account *Account, symbol string, stream pb.Admin_QuoteStreamServer) error {
	slog.Debug("client connected from quote stream", "symbol", symbol)

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

func (c *FinamClient) subscribeQuotes(symbol string, client *finam.Client, group *peers.Group[*marketdata.Quote]) {
	minSleep := time.Second
	maxSleep := time.Minute
	sleep := minSleep
LOOP:
	for {
		stream, err := client.SubscribeQuotes(context.Background(), []string{symbol})
		if err != nil {
			slog.Warn("cannot subscribe for quote stream", "symbol", symbol, "error", err.Error())
			if sleep > maxSleep {
				slog.Warn("max reconnects reached to quote stream", "symbol", symbol)
				break LOOP
			}
			time.Sleep(sleep)
			sleep = sleep * 2
			continue LOOP
		}

		sleep = minSleep

		for q, err := range stream {
			if err != nil {
				slog.Warn("quote stream disconnected", "symbol", symbol, "error", err.Error())
				break
			}
			if !group.Send(q) {
				slog.Warn("no subscribers for quote stream", "symbol", symbol)
				break LOOP
			}
		}
	}

	slog.Info("exit quote stream")
	group.Close()
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
