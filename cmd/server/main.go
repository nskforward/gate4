package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := app.NewApp()
	err := app.Start(ctx)
	if err != nil {
		slog.Error("cannot start app", "reason", err.Error())
		os.Exit(1)
	}
}
