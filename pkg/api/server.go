package api

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ernoaapa/eliot/pkg/model"

	"golang.org/x/net/context"

	"github.com/ernoaapa/eliot/pkg/api/mapping"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/api/stream"
	resolver "github.com/ernoaapa/eliot/pkg/node"
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Server implements the GRPC API for the eli
type Server struct {
	resolver *resolver.Resolver
	client   runtime.Client
	grpc     *grpc.Server
	listen   string
}

// Info is Node service Info implementation
func (s *Server) Info(context context.Context, req *node.InfoRequest) (*node.InfoResponse, error) {
	return &node.InfoResponse{
		Info: mapping.MapInfoToAPIModel(s.resolver.GetInfo()),
	}, nil
}

// Create is 'pods' service Create implementation
func (s *Server) Create(req *pods.CreatePodRequest, server pods.Pods_CreateServer) error {
	pod := mapping.MapPodToInternalModel(req.Pod)
	var (
		done       = make(chan struct{})
		progresses = []*progress.ImageFetch{}
	)
	defer close(done)

	if err := s.ensurePodNotExist(pod.Metadata.Namespace, pod.Metadata.Name); err != nil {
		return errors.Wrapf(err, "Cannot create pod [%s]", pod.Metadata.Name)
	}

	go func() {
		for {
			select {
			case <-done:
				// Send last update
				images := mapping.MapImageFetchProgressToAPIModel(progresses)

				if err := server.Send(&pods.CreatePodStreamResponse{Images: images}); err != nil {
					log.Warnf("Error while sending last create pod status back to client: %s", err)
				}
				return // End update loop
			case <-time.After(100 * time.Millisecond):
				images := mapping.MapImageFetchProgressToAPIModel(progresses)

				if err := server.Send(&pods.CreatePodStreamResponse{Images: images}); err != nil {
					log.Warnf("Error while sending create pod status back to client: %s", err)
				}
			}
		}
	}()

	for _, container := range pod.Spec.Containers {
		progress := progress.NewImageFetch(container.Name, container.Image)
		progresses = append(progresses, progress)

		if err := s.client.PullImage(pod.Metadata.Namespace, container.Image, progress); err != nil {
			progress.SetToFailed()
			return errors.Wrapf(err, "Failed to pull image [%s]", container.Image)
		}
		progress.AllDone()

		_, err := s.client.CreateContainer(pod, container)
		if err != nil {
			return errors.Wrapf(err, "Failed to create container [%s]", container.Name)
		}
		log.Debugf("Container [%s] created", container.Name)
	}

	return nil
}

func (s *Server) ensurePodNotExist(namespace, name string) error {
	_, err := s.client.GetPod(namespace, name)
	if err != nil {
		if runtime.IsNotFound(err) {
			return nil
		}
		return err
	}
	return fmt.Errorf("Pod [%s] in namespace [%s] already exist", name, namespace)
}

// Start is 'pods' service Start implementation
func (s *Server) Start(context context.Context, req *pods.StartPodRequest) (*pods.StartPodResponse, error) {
	pod, err := s.client.GetPod(req.Namespace, req.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find containers to start for pod [%s] in namespace [%s]", req.Name, req.Namespace)
	}

	iosets, err := buildContainerIOSets(pod.Metadata.Name, pod.Spec.Containers)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot start pod [%s], error while building IO sets for containers", req.Name)
	}

	statuses := []model.ContainerStatus{}
	for _, status := range pod.Status.ContainerStatuses {
		status, err := s.client.StartContainer(pod.Metadata.Namespace, status.ContainerID, *iosets[status.Name])
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to start container [%s]", status.Name)
		}
		log.Debugf("Container [%s] started", status.Name)
		statuses = append(statuses, status)
	}

	pod.Status.ContainerStatuses = statuses

	return &pods.StartPodResponse{
		Pod: mapping.MapPodToAPIModel(pod),
	}, nil
}

func buildContainerIOSets(podName string, containers []model.Container) (map[string]*runtime.IOSet, error) {
	iosets := map[string]*runtime.IOSet{}

	for _, container := range containers {
		ioset, err := runtime.NewIOSet(fmt.Sprintf("%s.%s", podName, container.Name))
		if err != nil {
			return nil, err
		}
		iosets[container.Name] = ioset
	}

	for _, container := range containers {
		ioset := iosets[container.Name]
		if container.Pipe != nil {
			target := iosets[container.Pipe.Stdout.Stdin.Name]
			if target == nil {
				return nil, fmt.Errorf("Invalid pipe definition, target container with name [%s] not found", container.Pipe.Stdout.Stdin.Name)
			}
			ioset.PipeStdoutTo(target)
		}
	}
	return iosets, nil
}

// Delete is 'pods' service Delete implementation
func (s *Server) Delete(context context.Context, req *pods.DeletePodRequest) (*pods.DeletePodResponse, error) {
	pod, err := s.client.GetPod(req.Namespace, req.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot fetch pod containers, cannot delete pod [%s]", req.Name)
	}

	statuses := []model.ContainerStatus{}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		status, err := s.client.StopContainer(req.Namespace, containerStatus.ContainerID)
		if err != nil {
			return nil, errors.Wrapf(err, "Error while stopping container [%s]", containerStatus.ContainerID)
		}
		statuses = append(statuses, status)
	}

	pod.Status.ContainerStatuses = statuses

	return &pods.DeletePodResponse{
		Pod: mapping.MapPodToAPIModel(pod),
	}, nil
}

// List is 'pods' service List implementation
func (s *Server) List(context context.Context, req *pods.ListPodsRequest) (*pods.ListPodsResponse, error) {
	p, err := s.client.GetPods(req.Namespace)
	if err != nil {
		return nil, err
	}
	return &pods.ListPodsResponse{
		Pods: mapping.MapPodsToAPIModel(p),
	}, nil
}

// Exec connects to process in container and streams stdout and stderr outputs to client
func (s *Server) Exec(server containers.Containers_ExecServer) error {
	md, ok := metadata.FromIncomingContext(server.Context())
	if !ok {
		return fmt.Errorf("Incoming Exec request don't have metadata. You must provide 'namespace', 'container', 'execid' and 'command' through metadata")
	}
	log.Debugf("Received metadata: %s", md)
	var (
		execID      = getMetadataValue(md, "execid")
		namespace   = getMetadataValue(md, "namespace")
		containerID = getMetadataValue(md, "container")
		args        = strings.Split(getMetadataValue(md, "args"), " ")
		tty         = false
	)
	tty, _ = strconv.ParseBool(getMetadataValue(md, "tty"))

	if namespace == "" {
		return fmt.Errorf("You must define 'namespace' metadata")
	}

	if containerID == "" {
		return fmt.Errorf("You must define 'container' metadata")
	}

	if len(args) == 0 {
		return fmt.Errorf("You must define 'args' metadata")
	}

	log.Debugf("Execute command [%s](tty: %t) in container [%s] in namespace [%s]", strings.Join(args, " "), tty, containerID, namespace)
	return s.client.Exec(
		namespace,
		containerID,
		execID,
		args,
		tty,
		runtime.AttachIO{
			Stdin:  stream.NewReader(server),
			Stdout: stream.NewWriter(server, false),
			Stderr: stream.NewWriter(server, true),
		},
	)
}

// Attach connects to process in container and streams stdout and stderr outputs to client
func (s *Server) Attach(server containers.Containers_AttachServer) error {
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

	log.Debugf("Attach to container [%s] in namespace [%s]", containerID, namespace)
	return s.client.Attach(
		namespace, containerID,
		runtime.AttachIO{
			Stdin:  stream.NewReader(server),
			Stdout: stream.NewWriter(server, false),
			Stderr: stream.NewWriter(server, true),
		},
	)
}

// Signal connects to process in container and send signal to the process
func (s *Server) Signal(cxt context.Context, req *containers.SignalRequest) (*containers.SignalResponse, error) {
	err := s.client.Signal(req.Namespace, req.ContainerID, syscall.Signal(req.Signal))
	if err != nil {
		return nil, err
	}
	return &containers.SignalResponse{}, nil
}

func getMetadataValue(md metadata.MD, key string) string {
	if val, ok := md[key]; ok {
		return val[0]
	}
	return ""
}

// NewServer creates new API server
func NewServer(listen string, client runtime.Client, resolver *resolver.Resolver) *Server {
	apiserver := &Server{
		resolver: resolver,
		client:   client,
		listen:   listen,
	}

	apiserver.grpc = grpc.NewServer()
	pods.RegisterPodsServer(apiserver.grpc, apiserver)
	containers.RegisterContainersServer(apiserver.grpc, apiserver)
	node.RegisterNodeServer(apiserver.grpc, apiserver)
	return apiserver
}

// Serve starts the server to serve GRPC server
func (s *Server) Serve() {
	log.Println("Start GRPC server...")
	lis, err := net.Listen("tcp", s.listen)
	if err != nil {
		log.Panicf("Failed to start API server to listen [%s]: %s", s.listen, err)
	}

	if err := s.grpc.Serve(lis); err != nil {
		log.Panicf("GRPC server stopped with error: %s", err)
	}
}

// Stop the GRPC server
func (s *Server) Stop() {
	log.Infof("Stop GRPC server...")
	s.grpc.Stop()
}
