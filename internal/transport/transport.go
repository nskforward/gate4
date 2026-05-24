package transport

import (
	"context"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type Transport struct {
	network string
	address string
	opts    []grpc.ServerOption
}

func NewTransport(network, address string, opt ...grpc.ServerOption) *Transport {
	return &Transport{
		network: network,
		address: address,
		opts:    opt,
	}
}

func (socket *Transport) Serve(ctx context.Context, gate4Server *Gate4Server) error {

	listener, err := net.Listen(socket.network, socket.address)
	if err != nil {
		return err
	}

	if socket.network == "unix" {
		defer os.Remove(socket.address)
	}

	grpcServer := grpc.NewServer(socket.opts...)
	pb.RegisterAdminServer(grpcServer, gate4Server)

	slog.Info("grpc server is ready to serve requests", "network", socket.network, "address", socket.address)

	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := grpcServer.Serve(listener)
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	select {
	case err := <-errorc:
		return err

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-ctx.Done():
			grpcServer.Stop()
		}
		return nil
	}
}
