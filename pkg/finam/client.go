package finam

import (
	"context"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *Conn
}

func NewClient(creds *Creds) (*Client, error) {
	conn, err := Connect(creds)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
	}
	return client, nil
}

func (client *Client) Close() error {
	client.cancel()
	return client.conn.Close()
}

func (client *Client) GetToken() (*Token, error) {
	return client.conn.GetToken()
}
