package server

import (
	"context"
	"fmt"
	"net"

	"github.com/nskforward/gate4/internal/api/grpc/server/handler"
	"github.com/nskforward/gate4/internal/api/grpc/server/interceptor"
	"github.com/nskforward/gate4/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TCPServer struct {
	addr      string
	transport *grpc.Server
}

func NewTCPServer(cfg config.Config, userHandler *handler.UserHandler, tokenHandler *handler.TokenHandler) (*TCPServer, error) {
	tlsConfig, err := NewTLSConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get tls config: %w", err)
	}

	s := &TCPServer{
		addr: cfg.GRPC.Addr,
		transport: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptor.Logging,
				interceptor.Recovery,
			),
			grpc.Creds(credentials.NewTLS(tlsConfig)),
		),
	}

	userHandler.Register(s.transport)
	tokenHandler.Register(s.transport)

	return s, nil
}

func (s *TCPServer) Register(handlers ...Registrable) {
	for _, h := range handlers {
		h.Register(s.transport)
	}
}

func (s *TCPServer) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	err = s.transport.Serve(l)
	if err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func (s *TCPServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.transport.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:

	case <-ctx.Done():
		s.transport.Stop()
	}
	return nil
}
