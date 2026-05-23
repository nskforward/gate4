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
	"github.com/nskforward/gate4/pkg/servers"
	google "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	tlsConfig, err := servers.MTLSConfig(cfg.CA.Cert, cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
	if err != nil {
		return err
	}

	var serverPool servers.Pool

	tcpServer := grpc.NewGate4Server(userStore, keychainStore, google.Creds(credentials.NewTLS(tlsConfig)))

	unixServer := grpc.NewGate4Server(userStore, keychainStore)

	serverPool.Add(func(poolCtx context.Context) error {
		return tcpServer.Run(poolCtx, "tcp", cfg.Server.TCPAddr)
	})

	serverPool.Add(func(poolCtx context.Context) error {
		return unixServer.Run(poolCtx, "unix", cfg.Server.UnixAddr)
	})

	return serverPool.Run(ctx)
}

func setLogLevel(level slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
