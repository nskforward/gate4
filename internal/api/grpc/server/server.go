package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/internal/api/grpc/server/interceptor"
	"github.com/nskforward/gate4/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	tcpAddr    string
	socketPath string
	tcpServer  *grpc.Server
	unixServer *grpc.Server
}

func NewServer(cfg config.Config,
	userHandler *UserHandler,
) (*Server, error) {

	tlsConfig, err := newTLSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot init tls config: %w", err)
	}

	interceptors := grpc.ChainUnaryInterceptor(
		interceptor.Logging,
		interceptor.Recovery,
	)

	tcpServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		interceptors,
	)
	unixServer := grpc.NewServer(interceptors)
	userHandler.Register(tcpServer, unixServer)
	return &Server{
		tcpAddr:    cfg.TCPAddr,
		socketPath: filepath.Join(os.TempDir(), "gate4.sock"),
		tcpServer:  tcpServer,
		unixServer: unixServer,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
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
