package api

import (
	"fmt"
	"net"

	"golang.org/x/net/context"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Server implements the GRPC API for the layery-cli
type Server struct {
	client runtime.Client
	grpc   *grpc.Server
	port   int
}

// List is pods List implementation
func (s *Server) List(context context.Context, req *pb.ListPodsRequest) (*pb.ListPodsResponse, error) {
	pods, err := s.client.GetContainersByPods(req.GetNamespace())
	if err != nil {
		return nil, err
	}
	return &pb.ListPodsResponse{
		Pods: mapPodsToApiModel(pods),
	}, nil
}

// NewServer creates new API server
func NewServer(port int, client runtime.Client) *Server {
	apiserver := &Server{
		client: client,
		port:   port,
	}

	apiserver.grpc = grpc.NewServer()
	pb.RegisterPodsServer(apiserver.grpc, apiserver)

	return apiserver
}

// Serve starts the server to serve GRPC server
func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		return errors.Wrapf(err, "Failed to start API server to listen port []", s.port)
	}
	return s.grpc.Serve(lis)
}
