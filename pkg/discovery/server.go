package discovery

import (
	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Server is zeroconf discovery server
type Server struct {
	Name   string
	Domain string
	Port   int
	server *zeroconf.Server
}

// NewServer creates new discovery server
func NewServer(name string, port int) *Server {
	return &Server{
		Name:   name,
		Domain: "local.",
		Port:   port,
	}
}

// Start server to be discoverable
func (s *Server) Start() error {
	log.Debugf("Exposing %s in port %d", s.Name, s.Port)
	server, err := zeroconf.Register(s.Name, ZeroConfServiceName, s.Domain, s.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to create zeroconf server")
	}
	s.server = server
	return nil
}

// Stop server to be discoverable
func (s *Server) Stop() {
	defer s.server.Shutdown()
}
