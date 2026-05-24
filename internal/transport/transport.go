package transport

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Listen(ctx context.Context, listener net.Listener, gate4Server *Gate4Server, tlsConfig *tls.Config) error {

	opts := []grpc.ServerOption{}
	mtls := "disabled"
	if tlsConfig != nil {
		mtls = "enabled"
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterGate4Server(grpcServer, gate4Server)

	slog.Info("grpc server is ready to serve requests", "network", listener.Addr().Network(), "address", listener.Addr().String(), "mtls", mtls)

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
			gate4Server.Close()
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
