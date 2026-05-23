package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/transport/grpc"
	"github.com/nskforward/gate4/internal/users"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	setLogLevel(slog.LevelDebug)

	slog.Info("start server")

	err := run(ctx)
	if err != nil {
		slog.Error("server stopped with error", "msg", err.Error())
	}
}

func run(ctx context.Context) error {
	cfg := config.Load()

	userStore, err := users.NewFileStorage()
	if err != nil {
		return err
	}

	keychainStore, err := keychain.NewStore(ctx, cfg)
	if err != nil {
		return err
	}

	server := grpc.NewGate4Server(userStore, keychainStore)

	return server.Run(ctx, "tcp", cfg.Server.TCPAddr)
}

func setLogLevel(level slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
