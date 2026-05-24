package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/servers"
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

	gate4Server, err := transport.NewGate4Server(cfg)
	if err != nil {
		return err
	}

	tlsConfig, err := servers.MTLSConfig(cfg.CA.Cert, cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
	if err != nil {
		return err
	}

	var serverPool servers.Pool

	serverPool.Add(func(poolCtx context.Context) error {
		socketPath := transport.UnixSocketPath()
		listener, err := net.Listen("unix", socketPath)
		if err != nil {
			return err
		}
		defer os.Remove(socketPath)
		return transport.Listen(poolCtx, listener, gate4Server, nil)
	})

	serverPool.Add(func(poolCtx context.Context) error {
		listener, err := net.Listen("tcp", cfg.Server.TCPAddr)
		if err != nil {
			return err
		}
		return transport.Listen(poolCtx, listener, gate4Server, tlsConfig)
	})

	return serverPool.Run(ctx)
}

func setLogLevel(level slog.Leveler) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
