package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/cli"
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
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	r := cli.NewRouter()
	routes(r, client)

	return r.Run(ctx)
}
