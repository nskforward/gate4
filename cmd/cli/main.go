package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/cli"
	"github.com/nskforward/gate4/internal/transport"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

func run(ctx context.Context) error {
	grpcClient, err := transport.NewGate4Client("unix", transport.UnixSocketPath())
	if err != nil {
		return err
	}
	defer grpcClient.Close()

	r := cli.NewRouter()
	r.Handle("help", cli.Help)
	r.Handle("cert create", cli.CreateCert(grpcClient))
	r.Handle("user list", cli.ListUsers(grpcClient))
	r.Handle("user create", cli.CreateUser(grpcClient))
	r.Handle("user delete", cli.DeleteUser(grpcClient))
	r.Handle("user block", cli.BlockUser(grpcClient))
	r.Handle("user edit", cli.UpdateUser(grpcClient))

	return r.Run(ctx)
}
