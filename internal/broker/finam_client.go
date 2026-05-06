package broker

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/internal/broker/types"
	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/peers"
	"github.com/nskforward/gate4/pkg/retries"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type FinamClient struct {
	clients       *finam.MultiClient
	quoteStream   *peers.PubSub[*marketdata.Quote]
	scheduleStore *ScheduleStore
	positionStore *PositionStore
}

func NewFinamClient() *FinamClient {
	return &FinamClient{
		clients:       finam.NewMultiClient(),
		quoteStream:   peers.NewPubSub[*marketdata.Quote](),
		scheduleStore: NewScheduleStore(),
		positionStore: NewPositionStore(),
	}
}

func (c *FinamClient) Positions(ctx context.Context, account *Account) ([]*pb.Position, error) {
	client, err := c.getClient(account)
	if err != nil {
		return nil, err
	}
	positions, err := c.positionStore.Get(ctx, client, account.ID)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.Position, 0, len(positions))
	for _, pos := range positions {
		result = append(result, &pb.Position{
			Symbol: pos.Symbol,
			Price:  pos.Price,
			Size:   pos.Size,
			Profit: pos.Profit,
		})
	}
	return result, nil
}

func (c *FinamClient) Schedule(ctx context.Context, account *Account, symbol string) ([]types.Session, types.Session, error) {
	client, err := c.getClient(account)
	if err != nil {
		return nil, types.Session{}, err
	}

	return c.scheduleStore.Sessions(ctx, client, symbol)
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
	return c.clients.Get(account.ID, account.Secret)
}

/*
func getPositions(in []*accounts.Position) []*pb.Position {
	items := make([]*pb.Position, 0, len(in))
	for _, item := range in {
		items = append(items, &pb.Position{
			Symbol:       item.Symbol,
			AveragePrice: item.AveragePrice.Value,
			Size:         item.Quantity.Value,
		})
	}
	return items
}
*/

func sessionTypeCast(in string) (types.SessionType, error) {
	sessionType, err := finam.NewSessionType(in)
	if err != nil {
		return types.SessionClosed, fmt.Errorf("cannot recognize the input session type '%s': %w", in, err)
	}

	switch sessionType {

	case finam.SessionClosed:
		return types.SessionClosed, nil

	case finam.SessionCoreTrading:
		return types.SessionMain, nil

	case finam.SessionOpeningAuction, finam.SessionEarlyTrading:
		return types.SessionPremarket, nil

	case finam.SessionClosingAuction, finam.SessionLateTrading:
		return types.SessionPostmarket, nil

	default:
		return types.SessionClosed, fmt.Errorf("cannot recognize the finam session type: %v", sessionType)
	}
}
