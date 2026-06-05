package server

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/nskforward/gate4/internal/api/grpc/server/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TCPServer struct {
	addr      string
	transport *grpc.Server
}

func NewTCPServer(addr string, tlsConfig *tls.Config) *TCPServer {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Recovery,
		),
	}

	if tlsConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	return &TCPServer{
		addr:      addr,
		transport: grpc.NewServer(opts...),
	}
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
