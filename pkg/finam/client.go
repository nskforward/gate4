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
func (c *Client) GetAsset(ctx context.Context, accountID, symbol string) (*types.AssetInfo, error) {
	reqCtx, cancel, err := c.authContextWithTimeout(ctx, 30*time.Second)
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
	result := convertAsset(resp)
	return &result, nil
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

// GetSchedule возвращает расписание торгов
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

	cache := make(map[string]*types.Quote)

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
		quotes := convertQuotes(resp.Quote)
		for _, q := range quotes {
			changed := changedQuote(cache, q)
			if changed != nil {
				if !send(*changed) {
					return nil
				}
			}
		}
	}
}

/*
func (c *Client) SubscribeQuotes(ctx context.Context, symbol string) *streams.Stream[types.Quote] {

	return c.quoteStreams.Subscribe(ctx, symbol, func(ctx context.Context, publish func(data types.Quote) bool) error {

		retry := retries.New(retries.Config{
			InitialDelay:  500 * time.Millisecond,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			MaxAttempts:   10,
			MaxJitter:     time.Second,
			OnBefore: func(attempt int) {
				slog.Debug("create a new outgoing finam quote stream", "symbol", symbol, "attempt", attempt)
			},
			OnAfter: func(err error) {
				slog.Debug("finam remote quote stream exited", "symbol", symbol, "error", err)
			},
		})

		return retry.Do(ctx, func(attempt int) error {

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

			cache := make(map[string]*types.Quote)

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
				quotes := convertQuotes(resp.Quote)
				for _, q := range quotes {
					changed := changedQuote(cache, q)
					if changed != nil {
						if !publish(*changed) {
							return nil
						}
					}
				}
			}
		})
	})
}
*/

/*
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
*/

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
