package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCTransport struct {
	tcpAddr    string
	socketPath string
	tcpServer  *grpc.Server
	unixServer *grpc.Server
}

func NewGRPCTransport(cfg config.Config,
	grpcUserHandler *handler.GRPCUserHandler,
) (*GRPCTransport, error) {

	tlsConfig, err := newTLSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot init tls config: %w", err)
	}
	tcpServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(),
	)
	unixServer := grpc.NewServer(grpc.ChainUnaryInterceptor())
	grpcUserHandler.Register(tcpServer, unixServer)
	return &GRPCTransport{
		tcpAddr:    cfg.TCPAddr,
		socketPath: filepath.Join(os.TempDir(), "gate4.sock"),
		tcpServer:  tcpServer,
		unixServer: unixServer,
	}, nil
}

func (s *GRPCTransport) Start(ctx context.Context) error {
	tcpListener, err := net.Listen("tcp", s.tcpAddr)
	if err != nil {
		return err
	}

	unixListener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	defer os.Remove(s.socketPath)

	errorc := make(chan error, 2)
	defer close(errorc)

	go serve(errorc, unixListener, s.unixServer)
	go serve(errorc, tcpListener, s.tcpServer)

	slog.Info("grpc server started")

	select {
	case <-ctx.Done():
		gracefulShutDown(s.unixServer, s.tcpServer)
		return nil

	case err := <-errorc:
		gracefulShutDown(s.unixServer, s.tcpServer)
		return err
	}
}

func serve(c chan error, l net.Listener, s *grpc.Server) {
	err := s.Serve(l)
	if err != nil && err != grpc.ErrServerStopped {
		c <- err
	}
}

func gracefulShutDown(servers ...*grpc.Server) {
	for _, s := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:

		case <-ctx.Done():
			s.Stop()
		}
	}
}
