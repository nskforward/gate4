package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/ssl"
	"golang.org/x/term"
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
		return errors.New("no commands provided, run `gate4 help` to see possible commands")
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
	case "account":
		return commandAccount(ctx, client, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func commandAccount(ctx context.Context, client *transport.AdminClient, args []string) error {
	if len(args) == 0 {
		return errors.New("broker command requres sub command")
	}
	subCommand := args[0]
	args = args[1:]
	switch subCommand {
	case "list":
		return commandAccountList(ctx, client)
	case "add":
		return commandAccountAdd(ctx, client)
	case "del":
		return commandAccountDelete(ctx, client, args)
	default:
		return fmt.Errorf("unknown account sub command: %s", subCommand)
	}
}

func commandAccountDelete(ctx context.Context, client *transport.AdminClient, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requres 1 argument")
	}
	return client.DeleteAccount(ctx, args[0])
}

func commandAccountAdd(ctx context.Context, client *transport.AdminClient) error {
	var brokerID, accountID, validUntil string
	fmt.Print("broker_id: ")
	fmt.Scanln(&brokerID)
	fmt.Print("account_id: ")
	fmt.Scanln(&accountID)
	fmt.Print("secret: ")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	secret, err := reader.ReadBytes('\r')
	if err != nil {
		panic(err)
	}
	if len(secret) > 0 && secret[len(secret)-1] == '\r' {
		secret = secret[:len(secret)-1]
	}
	term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Printf("[%d bytes]\n", len(secret))

	fmt.Print("valid until date (YYYY-MM-DD): ")
	fmt.Scanln(&validUntil)

	t, err := time.Parse("2006-01-02", validUntil)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	return client.AddAccount(ctx, brokerID, accountID, string(secret), t.Unix())
}

func commandAccountList(ctx context.Context, client *transport.AdminClient) error {
	items, err := client.ListAccounts(ctx)
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
	maxCommandLen := 18
	commandInfo := func(command, description, example string) {
		fmt.Printf(" %s %s %s\n %s  $ %s\n\n", command, strings.Repeat(" ", maxCommandLen-len(command)), description, strings.Repeat(" ", maxCommandLen), example)
	}
	printHeader("Help")
	commandInfo("help", "show help menu", "gate4 help")
	commandInfo("init-ca", "generate ssl CA key and cert and store to files", "gate4 init-ca path/to/key path/to/cert")
	commandInfo("issue-cert <name>", "generate and sign the robot ssl certicate", "gate4 issue-cert <name>")
	commandInfo("account list", "show accounts list", "gate4 account list")
	commandInfo("account add", "add a new account", "gate4 account add")
	commandInfo("account del <key>", "delete a account", "gate4 account del <key>")
	return nil
}

func printHeader(header string) {
	fmt.Println("=====================================================================")
	fmt.Println(" GATE4:", strings.ToUpper(header))
	fmt.Println("=====================================================================")
}
