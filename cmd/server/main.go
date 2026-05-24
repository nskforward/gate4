package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/servers"
	"google.golang.org/grpc"
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

	tlsConfig, err := servers.MTLSConfig(cfg.CA.Cert, cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
	if err != nil {
		return err
	}

	gate4Server, err := createGate4Server(ctx, cfg)
	if err != nil {
		return err
	}

	unixSocket := transport.NewTransport("unix", transport.UnixSocketPath())
	tcpSocket := transport.NewTransport("tcp", cfg.Server.TCPAddr, grpc.Creds(credentials.NewTLS(tlsConfig)))

	var serverPool servers.Pool
	serverPool.Add(func(poolCtx context.Context) error {
		return unixSocket.Serve(poolCtx, gate4Server)
	})
	serverPool.Add(func(poolCtx context.Context) error {
		return tcpSocket.Serve(poolCtx, gate4Server)
	})
	return serverPool.Run(ctx)
}

func createGate4Server(ctx context.Context, cfg config.Config) (*transport.Gate4Server, error) {
	userStore, err := users.NewFileStorage()
	if err != nil {
		return nil, err
	}
	keychainStore, err := keychain.NewStore(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return transport.NewGate4Server(userStore, keychainStore), nil
}

func setLogLevel(level slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
