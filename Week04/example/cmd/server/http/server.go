package http

import (
	"context"
	"log"
	"net/http"
)

// Server is a HTTP server.
type Server struct {
	srv *http.Server
}

// NewServer new a HTTP server.
func NewServer(mux *http.ServeMux) *Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Server{srv: srv}
}

// Start start the HTTP server.
func (s *Server) Start() error {
	log.Printf("[HTTP] Listening on: %s\n", s.srv.Addr)
	return s.srv.ListenAndServe()
}

// Stop shutdown the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Printf("[HTTP] Stopping")
	return s.srv.Shutdown(ctx)
}
