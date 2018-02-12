package profile

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

// NewServer creates new profiling server
func NewServer(addr string) *Server {
	return &Server{
		httpServer: &http.Server{Addr: addr},
	}
}

// Serve starts http server which serve the pperf
func (s *Server) Serve() {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Panicf("Failed to start profiling http server: %s", err)
	}
}

// Stop the http server
func (s *Server) Stop() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Panicf("Failed to stop profiling http server: %s", err)
	}
}
