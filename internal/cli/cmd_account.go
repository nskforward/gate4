package cli

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func (c *Client) cmdAccount(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("requres sub command")
	}
	command := args[0]
	args = args[1:]
	switch command {
	case "list":
		return c.cmdAccountList(ctx)
	case "add":
		return c.cmdAccountAdd(ctx)
	case "del":
		return c.cmdAccountDelete(ctx, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (c *Client) cmdAccountDelete(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requres 1 argument")
	}
	return c.adminClient.DeleteAccount(ctx, args[0])
}

func (c *Client) cmdAccountAdd(ctx context.Context) error {
	var brokerID, accountID, validUntil string
	fmt.Print("broker_id: ")
	fmt.Scanln(&brokerID)
	fmt.Print("account_id: ")
	fmt.Scanln(&accountID)

	secret := AskSecret("secret")

	fmt.Printf("[%d bytes]\n", len(secret))

	fmt.Print("valid until date (YYYY-MM-DD): ")
	fmt.Scanln(&validUntil)

	t, err := time.Parse("2006-01-02", validUntil)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	return c.adminClient.AddAccount(ctx, brokerID, accountID, string(secret), t.Unix())
}

func (c *Client) cmdAccountList(ctx context.Context) error {
	items, err := c.adminClient.ListAccounts(ctx)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		fmt.Println("no accounts")
		return nil
	}
	for i := range len(items) {
		fmt.Printf("%d. %s.%s (%s)\n", i+1, items[i].BrokerId, items[i].Id, time.Unix(items[i].ValidUntil, 0).Format("2006-01-02"))
	}
	return nil
}
