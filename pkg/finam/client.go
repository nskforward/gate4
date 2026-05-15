package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	APIHost = "api.finam.ru:443"
)

type Client struct {
	accountID         string
	markedDataService marketdata.MarketDataServiceClient
	orderService      orders.OrdersServiceClient
	assetsService     assets.AssetsServiceClient
	accountService    accounts.AccountsServiceClient
	authService       auth.AuthServiceClient
	tokenService      *tokenService
}

// NewClient создаёт новый клиент Finam
func NewClient(accountID, secret string) (*Client, error) {
	conn, err := connect(APIHost)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	authService := auth.NewAuthServiceClient(conn)
	accountService := accounts.NewAccountsServiceClient(conn)
	assetsService := assets.NewAssetsServiceClient(conn)
	orderService := orders.NewOrdersServiceClient(conn)
	markedDataService := marketdata.NewMarketDataServiceClient(conn)

	return &Client{
		accountID:         accountID,
		markedDataService: markedDataService,
		orderService:      orderService,
		assetsService:     assetsService,
		accountService:    accountService,
		authService:       authService,
		tokenService:      newTokenService(authService, secret),
	}, nil
}

// GetAccountInfo возвращает информацию о счёте
func (c *Client) GetAccount(ctx context.Context) (*types.AccountInfo, error) {
	reqCtx, cancel, err := c.authContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()
	resp, err := c.accountService.GetAccount(reqCtx, &accounts.GetAccountRequest{
		AccountId: c.accountID,
	})
	if err != nil {
		return nil, fmt.Errorf("get account info failed: %w", err)
	}
	return &types.AccountInfo{
		AccountID: resp.AccountId,
		Positions: convertPositions(resp.Positions),
	}, nil
}

// GetAssetInfo возвращает информацию об активе
func (c *Client) GetAsset(ctx context.Context, symbol string) (types.AssetInfo, error) {
	reqCtx, cancel, err := c.authContextWithTimeout(ctx, 30*time.Second)
	if err != nil {
		return types.AssetInfo{}, err
	}
	defer cancel()
	resp, err := c.assetsService.GetAsset(reqCtx, &assets.GetAssetRequest{
		AccountId: c.accountID,
		Symbol:    symbol,
	})
	if err != nil {
		return types.AssetInfo{}, fmt.Errorf("get asset info failed: %w", err)
	}
	return convertAsset(resp), nil
}

func (c *Client) GetSchedule(ctx context.Context, symbol string) ([]types.Session, error) {
	reqCtx, cancel, err := c.authContextWithTimeout(ctx, 30*time.Second)
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
	return convertSessions(resp.Sessions), nil
}

/*
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
*/

/*
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
*/

/*
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
*/

func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) bool) error {
	reqCtx, err := c.authContext(ctx)
	if err != nil {
		return err
	}

	finamStream, err := c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: []string{symbol},
	})
	if err != nil {
		return err
	}

	cache := NewQuoteCache()

	for {
		resp, err := finamStream.Recv()
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Canceled {
					return nil
				}
			}
			return err
		}
		for _, q := range resp.Quote {
			result := cache.Get(q)
			if result != nil {
				if !send(*result) {
					return nil
				}
			}
		}
	}
}

func (c *Client) SubscribeAccountTrades(ctx context.Context, send func(types.AccountTrade) bool) error {
	reqCtx, err := c.authContext(ctx)
	if err != nil {
		return err
	}
	finamStream, err := c.orderService.SubscribeTrades(reqCtx, &orders.SubscribeTradesRequest{
		AccountId: c.accountID,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := finamStream.Recv()
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Canceled {
					return nil
				}
			}
			return err
		}
		for _, trade := range resp.Trades {
			if !send(convertAccountTrade(trade)) {
				return nil
			}
		}
	}
}

// ContextWithTimeout создаёт контекст с таймаутом
func (c *Client) authContextWithTimeout(ctx context.Context, timeout time.Duration) (reqCtx context.Context, cancel context.CancelFunc, err error) {
	reqCtx, cancel = context.WithTimeout(ctx, timeout)
	reqCtx, err = c.authContext(reqCtx)
	return
}

// Context добавляет токен авторизации к контексту
func (c *Client) authContext(ctx context.Context) (reqCtx context.Context, err error) {
	return c.tokenService.Context(ctx)
}
