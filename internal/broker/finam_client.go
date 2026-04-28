package broker

import (
	"context"
	"sync"

	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
)

type FinamClient struct {
	clients map[string]*finam.Client
	mx      sync.Mutex
}

func NewFinamClient() *FinamClient {
	return &FinamClient{
		clients: make(map[string]*finam.Client),
	}
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
