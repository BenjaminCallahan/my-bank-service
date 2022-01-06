package api

import (
	"context"
	"net/http"
	"time"
)

// Server represent entity of web server
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new entity of Server
func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

// Run starting web server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown shutting down a server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
