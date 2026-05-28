package app

import (
	"context"
	"crypto/tls"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nskforward/gate4/internal/api"
	"google.golang.org/grpc"
)

type JointServer struct {
	unixServer *UnixServer
	tcpServer  *TCPServer
}

func NewJointServer(apiServer *api.Server, tlsConfig *tls.Config, tcpAddr string) *JointServer {
	return &JointServer{
		unixServer: NewUnixServer(apiServer),
		tcpServer:  NewTCPServer(apiServer, tlsConfig, tcpAddr),
	}
}

func (s *JointServer) Close() {
	s.unixServer.Close()
	s.tcpServer.Close()
}

func (s *JointServer) Serve() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errorc := make(chan error, 2)
	defer close(errorc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := s.unixServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := s.tcpServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	select {
	case err := <-errorc:
		s.Close()
		return err

	case <-ctx.Done():
		s.Close()
		return nil
	}
}
