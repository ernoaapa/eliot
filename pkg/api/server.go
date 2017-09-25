package api

import (
	"net"

	"golang.org/x/net/context"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/api/stream"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	containersByPods, err := s.client.GetContainers(req.GetNamespace())
	if err != nil {
		return nil, err
	}
	return &pb.ListPodsResponse{
		Pods: mapPodsToAPIModel(req.GetNamespace(), containersByPods),
	}, nil
}

// Logs returns container logs
func (s *Server) Logs(req *pb.GetLogsRequest, resp pb.Pods_LogsServer) error {
	log.Debugf("Get logs for container [%s] in namespace [%s]", req.GetContainerID(), req.GetNamespace())
	return s.client.GetLogs(
		req.GetNamespace(), req.GetContainerID(),
		&stream.EmptyStdin{},
		stream.NewLogsWriter(resp, pb.GetLogsResponse_STDOUT),
		stream.NewLogsWriter(resp, pb.GetLogsResponse_STDERR),
	)
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
