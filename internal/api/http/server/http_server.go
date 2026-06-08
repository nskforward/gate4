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

func NewHTTPServer(cfg config.Config, mux *http.ServeMux) (*HTTPServer, error) {
	tlsConfig, err := NewTLSConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		transport: &http.Server{
			Addr:      cfg.HTTP.Addr,
			Handler:   mux,
			TLSConfig: tlsConfig,
		},
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context) error {
	err := s.transport.ListenAndServeTLS("", "")
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
