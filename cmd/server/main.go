package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/store"
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

	accountStore, err := store.NewAccountStore(store.NewAccountFileProvider(cfg.StoreDir))
	if err != nil {
		return err
	}

	adminServer := transport.NewAdminServer(cfg, logger, accountStore)
	gatewayServer := transport.NewGatewayServer(cfg, logger, accountStore)

	return race.Run(ctx, gatewayServer.Run, adminServer.Run)
}

func initLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	return logger
}
