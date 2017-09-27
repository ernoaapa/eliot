package api

import (
	"net"

	"golang.org/x/net/context"

	"github.com/ernoaapa/can/pkg/api/mapping"
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

// Create is 'pods' service Create implementation
func (s *Server) Create(context context.Context, req *pb.CreatePodRequest) (*pb.CreatePodResponse, error) {
	pod := mapping.MapPodToInternalModel(req.Pod)

	for _, container := range pod.Spec.Containers {
		if err := s.client.CreateContainer(pod, container); err != nil {
			return nil, errors.Wrap(err, "Failed to create container")
		}
		if err := s.client.StartContainer(pod.Metadata.Namespace, container.Name); err != nil {
			return nil, errors.Wrap(err, "Failed to start container")
		}
	}

	containers, err := s.client.GetContainers(pod.Metadata.Namespace, pod.Metadata.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch created containers info")
	}

	return &pb.CreatePodResponse{
		Pod: mapping.MapPodToAPIModel(pod.Metadata.Namespace, pod.Metadata.Name, containers),
	}, nil
}

// List is 'pods' service List implementation
func (s *Server) List(context context.Context, req *pb.ListPodsRequest) (*pb.ListPodsResponse, error) {
	containersByPods, err := s.client.GetAllContainers(req.Namespace)
	if err != nil {
		return nil, err
	}
	return &pb.ListPodsResponse{
		Pods: mapping.MapPodsToAPIModel(req.Namespace, containersByPods),
	}, nil
}

// Logs returns container logs
func (s *Server) Logs(req *pb.GetLogsRequest, server pb.Pods_LogsServer) error {
	log.Debugf("Get logs for container [%s] in namespace [%s]", req.GetContainerID(), req.Namespace)
	return s.client.GetLogs(
		req.Namespace, req.GetContainerID(),
		runtime.AttachIO{
			Stdin:  &stream.EmptyStdin{},
			Stdout: stream.NewWriter(server, false),
			Stderr: stream.NewWriter(server, true),
		},
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
