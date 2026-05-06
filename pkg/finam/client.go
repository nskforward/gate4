package finam

import (
	"context"
	"fmt"
	"io"
	"iter"
	"time"

	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"github.com/nskforward/gate4/pkg/peers"
)

const (
	APIHost = "api.finam.ru:443"
)

type Client struct {
	markedDataService marketdata.MarketDataServiceClient
	orderService      orders.OrdersServiceClient
	assetsService     assets.AssetsServiceClient
	accountService    accounts.AccountsServiceClient
	authService       auth.AuthServiceClient
	tokenService      *tokenService
	quoteStreams      map[string]*peers.Group[*marketdata.Quote]
}

// NewClient создаёт новый клиент Finam
func NewClient(addr, secret string) (*Client, error) {
	conn, err := connect(addr)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	authService := auth.NewAuthServiceClient(conn)
	accountService := accounts.NewAccountsServiceClient(conn)
	assetsService := assets.NewAssetsServiceClient(conn)
	orderService := orders.NewOrdersServiceClient(conn)
	markedDataService := marketdata.NewMarketDataServiceClient(conn)

	return &Client{
		markedDataService: markedDataService,
		orderService:      orderService,
		assetsService:     assetsService,
		accountService:    accountService,
		authService:       authService,
		tokenService:      newTokenService(authService, secret),
	}, nil
}

func (c *Client) UpdateSecret(secret string) error {
	return c.tokenService.updateSecret(secret)
}

// GetAccountInfo возвращает информацию о счёте
func (c *Client) GetAccountInfo(ctx context.Context, accountID string) (*accounts.GetAccountResponse, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	resp, err := c.accountService.GetAccount(reqCtx, &accounts.GetAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, fmt.Errorf("get account info failed: %w", err)
	}
	return resp, nil
}

// GetAssetInfo возвращает информацию об активе
func (c *Client) GetAssetInfo(ctx context.Context, accountID, symbol string) (*assets.GetAssetResponse, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	resp, err := c.assetsService.GetAsset(reqCtx, &assets.GetAssetRequest{
		AccountId: accountID,
		Symbol:    symbol,
	})
	if err != nil {
		return nil, fmt.Errorf("get asset info failed: %w", err)
	}
	return resp, nil
}

// CancelOrder отменяет ордер
func (c *Client) CancelOrder(ctx context.Context, accountID, orderID string) (*orders.OrderState, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	state, err := c.orderService.CancelOrder(reqCtx, &orders.CancelOrderRequest{
		AccountId: accountID,
		OrderId:   orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("cancel order failed: %w", err)
	}
	return state, nil
}

// GetOrders возвращает список ордеров
func (c *Client) GetOrders(ctx context.Context, accountID string) ([]*orders.OrderState, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	state, err := c.orderService.GetOrders(reqCtx, &orders.OrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, fmt.Errorf("get orders failed: %w", err)
	}
	return state.GetOrders(), nil
}

// GetSchedule возвращает расписание торгов
func (c *Client) GetSchedule(ctx context.Context, symbol string) ([]*assets.ScheduleResponse_Sessions, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	req := &assets.ScheduleRequest{
		Symbol: symbol,
	}
	resp, err := c.assetsService.Schedule(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("get schedule failed: %w", err)
	}
	return resp.Sessions, nil
}

// PlaceOrder размещает ордер
func (c *Client) PlaceOrder(ctx context.Context, order *orders.Order) (*orders.OrderState, error) {
	reqCtx, cancel, err := c.ContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	state, err := c.orderService.PlaceOrder(reqCtx, order)
	if err != nil {
		return nil, fmt.Errorf("place order failed: %w", err)
	}
	return state, nil
}

// SubscribeQuotes подписывается на котировки
func (c *Client) SubscribeQuotes(ctx context.Context, symbols []string) (iter.Seq2[*marketdata.Quote, error], error) {
	reqCtx, err := c.Context(ctx)
	if err != nil {
		return nil, err
	}
	stream, err := c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: symbols,
	})
	if err != nil {
		return nil, err
	}
	return func(yield func(*marketdata.Quote, error) bool) {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(nil, err)
				return
			}
			for _, order := range resp.Quote {
				if !yield(order, nil) {
					return
				}
			}
		}
	}, nil
}

// SubscribeAccountTrades подписывается на сделки аккаунта
func (c *Client) SubscribeAccountTrades(ctx context.Context, accountID string) (iter.Seq2[*v1.AccountTrade, error], error) {
	reqCtx, err := c.Context(ctx)
	if err != nil {
		return nil, err
	}
	stream, err := c.orderService.SubscribeTrades(reqCtx, &orders.SubscribeTradesRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}
	return func(yield func(*v1.AccountTrade, error) bool) {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(nil, err)
				return
			}
			for _, trade := range resp.Trades {
				if !yield(trade, nil) {
					return
				}
			}
		}
	}, nil
}

// ContextWithTimeout создаёт контекст с таймаутом
func (c *Client) ContextWithTimeout(ctx context.Context, timeout time.Duration) (reqCtx context.Context, cancel context.CancelFunc, err error) {
	reqCtx, cancel = context.WithTimeout(ctx, timeout)
	reqCtx, err = c.Context(reqCtx)
	return
}

// Context добавляет токен авторизации к контексту
func (c *Client) Context(ctx context.Context) (reqCtx context.Context, err error) {
	return c.tokenService.Context(ctx)
}
