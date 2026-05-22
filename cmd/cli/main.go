package main

import (
	"context"
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
	r.Handle("users", cli.ListUsers(adminClient))
	r.Handle("user add", cli.AddUser(adminClient))
	r.Handle("user del", cli.DeleteUser(adminClient))

	return r.Run(ctx)
}
