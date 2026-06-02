package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/cmd/client/handler"
	"github.com/nskforward/gate4/internal/api/grpc/client"
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

	r := handler.NewRouter()
	routes(r, c)

	return r.Run(ctx)
}
