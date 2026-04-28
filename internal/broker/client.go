package broker

import (
	"context"

	"github.com/nskforward/gate4/pkg/pb"
)

type Client interface {
	GetAccountInfo(ctx context.Context, account *Account) (*pb.AccountResponse, error)
}
