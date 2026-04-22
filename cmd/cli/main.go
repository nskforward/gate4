package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/ssl"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]:", err)
	}
}

func run(ctx context.Context) error {
	args := os.Args
	if len(args) < 2 {
		return errors.New("no commands provided")
	}
	client, err := transport.NewAdminClient(":4001")
	if err != nil {
		return fmt.Errorf("cannot connect to server: %w", err)
	}
	defer client.Close()
	return execCommand(ctx, client, args[1], args[2:])
}

func execCommand(ctx context.Context, client *transport.AdminClient, command string, args []string) error {
	switch command {
	case "help":
		return commandHelp(ctx)
	case "init-ca":
		return commandInitCA(ctx, args)
	case "broker":
		return commandBroker(ctx, client, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func commandBroker(ctx context.Context, client *transport.AdminClient, args []string) error {
	if len(args) == 0 {
		return errors.New("broker command requres sub command")
	}
	subCommand := args[0]
	args = args[1:]
	switch subCommand {
	case "list":
		return commandBrokerList(ctx, client)
	default:
		return fmt.Errorf("unknown broker sub command: %s", subCommand)
	}
}

func commandBrokerList(ctx context.Context, client *transport.AdminClient) error {
	items, err := client.ListBrokers(ctx)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		fmt.Println("no brokers")
		return nil
	}
	for i := range len(items) {
		fmt.Printf("%d. %s (%s)\n", i+1, items[i].Id, items[i].Address)
	}
	return nil
}

func commandInitCA(_ context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("requires 2 arguments")
	}
	err := ssl.GenCA(args[0], args[1])
	if err != nil {
		return fmt.Errorf("cannot generate CA: %w", err)
	}
	fmt.Println("success")
	return nil
}

func commandHelp(context.Context) error {
	maxCommandLen := 12
	commandInfo := func(command, description, example string) {
		fmt.Printf(" %s %s %s\n %s  $ %s\n\n", command, strings.Repeat(" ", maxCommandLen-len(command)), description, strings.Repeat(" ", maxCommandLen), example)
	}
	printHeader("Command Help")
	commandInfo("help", "show help menu", "gate4 help")
	commandInfo("init-ca", "generate ssl CA key and cert and store to files", "gate4 init-ca path/to/key path/to/cert")
	commandInfo("broker list", "show brokers list", "gate4 broker list")
	return nil
}

func printHeader(header string) {
	fmt.Println("============================================================")
	fmt.Println(" GATE4:", header)
	fmt.Println("============================================================")
}
