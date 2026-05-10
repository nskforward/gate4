package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/race"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("start application")

	err := run(context.Background())
	if err != nil {
		slog.Error("application exited with error", "error", err.Error())
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}
	cfg.LogParams()

	b, err := broker.NewBroker()
	if err != nil {
		return err
	}

	adminServer, err := transport.NewAdminServer(cfg, b)
	if err != nil {
		return err
	}
	gatewayServer, err := transport.NewGatewayServer(cfg, b)
	if err != nil {
		return err
	}

	return race.Run(ctx, gatewayServer.Run, adminServer.Run)
}
