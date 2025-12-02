package server

import (
	"context"
	"net"
	"net/http"

	"greddit/internal/infra/http/routing"
)

// Server represents an HTTP server.
type Server struct {
	config  Config
	handler http.Handler
	params  routing.RouterParams
	ln      net.Listener
}

// New creates a new HTTP server.
func New(p routing.RouterParams, opts ...Option) (*Server, error) {
	config := defaultConfig()

	for _, opt := range opts {
		opt(&config)
	}

	s := &Server{
		config: config,
		params: p,
	}

	if s.config.ln != nil {
		s.config.addr = s.config.ln.Addr().String()
	}

	mux := NewHandler(p)
	s.handler = LoggerMiddleware(mux, p.Logger)

	s.handler = corsMiddleware(s.handler, config.allowedOrigins)

	s.handler = TracerMiddleware(s.handler)

	return s, nil
}

func (s Server) Addr() string {
	return s.config.addr
}

// Start starts the HTTP server.
// It blocks until the server stops.
// Also performs a graceful shutdown according to the shutdown timeout.
func (s Server) Start(ctx context.Context, readyCh ...chan struct{}) error {
	logger := s.params.Logger

	server := http.Server{
		Addr:    s.Addr(),
		Handler: s.handler,
	}

	ch := make(chan error, 1)
	go func() {
		logger.Info("HTTP server started",
			"address", server.Addr)
		if s.config.ln == nil {
			ch <- server.ListenAndServe()
		} else {
			ch <- server.Serve(s.config.ln)
		}
	}()
	for _, ch := range readyCh {
		ch <- struct{}{}
	}

	defer logger.Info("HTTP server shut down completed")

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.shutdownTimeout)
		defer cancel()
		logger.Warn("HTTP server received context cancellation, attempting shutdown")

		return server.Shutdown(shutdownCtx)

	case err := <-ch:
		logger.Error("HTTP server received error", "error", err)

		return err
	}
}
