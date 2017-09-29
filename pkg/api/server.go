package api

import (
	"fmt"
	"net"

	"golang.org/x/net/context"

	"github.com/ernoaapa/can/pkg/api/mapping"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/api/stream"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
		log.Debugf("Container [%s] created, will start it", container.Name)
		if err := s.client.StartContainer(pod.Metadata.Namespace, container.Name, container.Tty); err != nil {
			return nil, errors.Wrap(err, "Failed to start container")
		}
		log.Debugf("Container [%s] started", container.Name)
	}

	containers, err := s.client.GetContainers(pod.Metadata.Namespace, pod.Metadata.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch created containers info")
	}

	return &pb.CreatePodResponse{
		Pod: mapping.MapPodToAPIModel(pod.Metadata.Namespace, pod.Metadata.Name, containers),
	}, nil
}

// Delete is 'pods' service Delete implementation
func (s *Server) Delete(context context.Context, req *pb.DeletePodRequest) (*pb.DeletePodResponse, error) {
	containers, err := s.client.GetContainers(req.Namespace, req.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot fetch pod containers, cannot delete pod [%s]", req.Name)
	}

	for _, container := range containers {
		err := s.client.StopContainer(req.Namespace, container.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "Error while stopping container [%s]", container.Name)
		}
	}
	return &pb.DeletePodResponse{
		Pod: mapping.MapPodToAPIModel(req.Namespace, req.Name, containers),
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

// Attach connects to process in container and and streams stdout and stderr outputs to client
func (s *Server) Attach(server pb.Pods_AttachServer) error {
	md, ok := metadata.FromIncomingContext(server.Context())
	if !ok {
		return fmt.Errorf("Incoming attach request don't have metadata. You must provide 'Namespace' and 'ContainerID' through metadata")
	}
	log.Debugf("Received metadata: %s", md)
	var (
		namespace   = getMetadataValue(md, "namespace")
		containerID = getMetadataValue(md, "container")
	)

	if namespace == "" {
		return fmt.Errorf("You must define 'namespace' metadata")
	}

	if containerID == "" {
		return fmt.Errorf("You must define 'container' metadata")
	}

	log.Debugf("Get logs for container [%s] in namespace [%s]", containerID, namespace)
	err := s.client.Attach(
		namespace, containerID,
		runtime.AttachIO{
			Stdin:  stream.NewReader(server),
			Stdout: stream.NewWriter(server, false),
			Stderr: stream.NewWriter(server, true),
		},
	)

	log.Printf("client.Attach ended with error: %s", err)
	return err
}

func getMetadataValue(md metadata.MD, key string) string {
	if val, ok := md[key]; ok {
		return val[0]
	}
	return ""
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
