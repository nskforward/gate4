package broker

import (
	"context"

	"github.com/nskforward/gate4/internal/broker/types"
	"github.com/nskforward/gate4/pkg/pb"
)

type Client interface {
	GetAccountInfo(ctx context.Context, account *Account) (*pb.AccountResponse, error)
	SubscribeQuotes(account *Account, symbol string, stream pb.Admin_QuoteStreamServer) error
	Schedule(ctx context.Context, account *Account, symbol string) ([]types.Session, types.Session, error)
}
