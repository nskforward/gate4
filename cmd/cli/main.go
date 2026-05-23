package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/transport/cli"
	"github.com/nskforward/gate4/internal/transport/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := run(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println()
			return
		}
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

func run(ctx context.Context) error {
	cfg := config.Load()
	adminClient, err := grpc.NewAdminClient(cfg.Admin.ListenAddr)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	r := cli.NewRouter()
	r.Handle("help", cli.Help)
	r.Handle("show users", cli.ListUsers(adminClient))
	r.Handle("create user", cli.AddUser(adminClient))
	r.Handle("delete user", cli.DeleteUser(adminClient))
	r.Handle("block user", cli.BlockUser(adminClient))
	r.Handle("edit user", cli.UpdateUser(adminClient))

	return r.Run(ctx)
}
