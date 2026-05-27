package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/app/client"
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
	app, err := client.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	r := cli.NewRouter()
	routes(r, app)

	return r.Run(ctx)
}
