package finam

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/pkg/types"
)

type Client struct {
}

func newClient(accountID, secret string) (*Client, error) {
	return nil, errors.New("finam client constructor not implemented")
}

func (c *Client) SubscribeQuotes(context.Context, func(types.Quote) error) error {
	return errors.New("finam: not implemented")
}
