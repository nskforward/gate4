package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/nskforward/gate4/internal/config"
)

type HTTPServer struct {
	transport *http.Server
}

func NewHTTPServer(cfg config.Config) *HTTPServer {
	return &HTTPServer{
		transport: &http.Server{
			Addr: cfg.HTTP.Addr,
		},
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	err := s.transport.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.transport.Shutdown(context.Background())
		close(stopped)
	}()

	select {
	case <-stopped:

	case <-ctx.Done():
		s.transport.Close()
	}
	return nil
}
