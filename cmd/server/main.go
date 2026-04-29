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
	logger := initLogger()
	logger.Info("start application")
	err := run(context.Background(), logger)
	if err != nil {
		logger.Error("application exited with error", "error", err.Error())
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}
	cfg.LogParams(logger)

	b, err := broker.NewBroker()
	if err != nil {
		return err
	}

	adminServer, err := transport.NewAdminServer(cfg, logger, b)
	if err != nil {
		return err
	}
	gatewayServer := transport.NewGatewayServer(cfg, logger, b)

	return race.Run(ctx, gatewayServer.Run, adminServer.Run)
}

func initLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	return logger
}
