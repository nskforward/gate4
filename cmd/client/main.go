package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/router"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	r := router.NewRouter("GATE 4 CLI client v0.0.1")
	routes(r, c)

	return r.Run(ctx, os.Args)
}
