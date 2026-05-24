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
	routes(r, grpcClient)

	return r.Run(ctx)
}
