package api

import (
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
	listen string
}

// List is 'pods' service List implementation
func (s *Server) List(context context.Context, req *pb.ListPodsRequest) (*pb.ListPodsResponse, error) {
	pods, err := s.client.GetContainers(req.GetNamespace())
	if err != nil {
		return nil, err
	}
	return &pb.ListPodsResponse{
		Pods: mapPodsToApiModel(pods),
	}, nil
}

// NewServer creates new API server
func NewServer(listen string, client runtime.Client) *Server {
	apiserver := &Server{
		client: client,
		listen: listen,
	}

	apiserver.grpc = grpc.NewServer()
	pb.RegisterPodsServer(apiserver.grpc, apiserver)

	return apiserver
}

// Serve starts the server to serve GRPC server
func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.listen)
	if err != nil {
		return errors.Wrapf(err, "Failed to start API server to listen [%s]", s.listen)
	}
	return s.grpc.Serve(lis)
}
