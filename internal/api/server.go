package api

import (
	"context"
	"log/slog"
	"sync"
	"time"

	grpcserver "github.com/nskforward/gate4/internal/api/grpc/server"
	httpserver "github.com/nskforward/gate4/internal/api/http/server"
)

type ServerInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Server struct {
	children []ServerInterface
}

func NewServer(unixServer *grpcserver.UnixServer, tcpServer *grpcserver.TCPServer, httpServer *httpserver.HTTPServer) *Server {
	return &Server{
		children: []ServerInterface{unixServer, tcpServer, httpServer},
	}
}

func (s *Server) Start(ctx context.Context) error {
	errorc := make(chan error, len(s.children))
	instances := 0
	for _, child := range s.children {
		instances++
		go func(ser ServerInterface) {
			err := ser.Start(ctx)
			if err != nil {
				errorc <- err
			}
		}(child)
	}

	go func() {
		defer close(errorc)
		select {
		case <-ctx.Done():
			stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			s.Stop(stopCtx)

		case err := <-errorc:
			slog.Error("server instance exited with error", "reason", err.Error())
			stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			s.Stop(stopCtx)
			errorc <- err
		}
	}()

	select {
	case err := <-errorc:
		return err
	case <-time.After(200 * time.Millisecond):
		slog.Info("server started", "instances", instances)
	}

	return <-errorc
}

func (s *Server) Stop(ctx context.Context) error {
	errorc := make(chan error, len(s.children))

	var wg sync.WaitGroup

	for _, child := range s.children {
		wg.Add(1)
		go func(ser ServerInterface) {
			defer wg.Done()
			err := ser.Stop(ctx)
			if err != nil {
				errorc <- err
			}
		}(child)
	}

	wg.Wait()

	close(errorc)

	return <-errorc
}
