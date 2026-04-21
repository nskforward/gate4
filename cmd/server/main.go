package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/transport"
)

func main() {
	logger := initLogger()
	logger.Info("start application")
	err := run(context.Background())
	if err != nil {
		logger.Error("application exited with error", "error", err.Error())
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}
	server := transport.NewServer(cfg.Addr)
	return server.Run(ctx)
}

func initLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	return logger
}
