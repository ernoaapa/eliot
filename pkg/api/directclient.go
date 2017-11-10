package api

import (
	"fmt"
	"io"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ernoaapa/elliot/pkg/api/mapping"
	containers "github.com/ernoaapa/elliot/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/elliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/elliot/pkg/api/stream"
	"github.com/ernoaapa/elliot/pkg/config"
	"github.com/ernoaapa/elliot/pkg/progress"
)

// DirectClient connects directly to device RPC API
type DirectClient struct {
	Namespace string
	Endpoint  config.Endpoint
	ctx       context.Context
}

// NewDirectClient creates new RPC server client
func NewDirectClient(namespace string, endpoint config.Endpoint) *DirectClient {
	return &DirectClient{
		namespace,
		endpoint,
		context.Background(),
	}
}

// GetPods calls server and fetches all pods information
func (c *DirectClient) GetPods() ([]*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	resp, err := client.List(c.ctx, &pods.ListPodsRequest{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPods(), nil
}

// GetPod return Pod by name
func (c *DirectClient) GetPod(podName string) (*pods.Pod, error) {
	pods, err := c.GetPods()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		if pod.Metadata.Name == podName {
			return pod, nil
		}
	}
	return nil, fmt.Errorf("Pod with name [%s] not found", podName)
}

// CreatePod creates new pod to the target server
func (c *DirectClient) CreatePod(status chan<- []*progress.ImageFetch, pod *pods.Pod, opts ...PodOpts) error {
	for _, o := range opts {
		err := o(pod)
		if err != nil {
			return err
		}
	}

	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	stream, err := client.Create(c.ctx, &pods.CreatePodRequest{
		Pod: pod,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			err = stream.CloseSend()
			return err
		}
		if err != nil {
			return err
		}

		status <- mapping.MapAPIModelToImageFetchProgress(resp.Images)
	}
}

// StartPod creates new pod to the target server
func (c *DirectClient) StartPod(name string) (*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	resp, err := client.Start(c.ctx, &pods.StartPodRequest{
		Namespace: c.Namespace,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPod(), nil
}

// DeletePod creates new pod to the target server
func (c *DirectClient) DeletePod(pod *pods.Pod) (*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)

	resp, err := client.Delete(c.ctx, &pods.DeletePodRequest{
		Namespace: pod.Metadata.Namespace,
		Name:      pod.Metadata.Name,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetPod(), nil
}

// Attach calls server and fetches pod logs
func (c *DirectClient) Attach(containerID string, attachIO AttachIO, hooks ...AttachHooks) (err error) {
	done := make(chan struct{})
	errc := make(chan error)

	md := metadata.Pairs(
		"namespace", c.Namespace,
		"container", containerID,
	)
	ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(c.ctx, md))
	defer cancel()

	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := containers.NewContainersClient(conn)
	log.Debugf("Open stream connection to server to get logs")
	s, err := client.Attach(ctx)
	if err != nil {
		return err
	}

	go func() {
		errc <- stream.PipeStdout(s, attachIO.Stdout, attachIO.Stderr)
	}()

	if attachIO.Stdin != nil {
		go func() {
			errc <- stream.PipeStdin(s, attachIO.Stdin)
		}()
	}

	for _, hook := range hooks {
		go hook(c.Endpoint, done)
	}

	for {
		err := <-errc
		close(done)
		return err
	}
}

// Signal sends kill signal to container process
func (c *DirectClient) Signal(containerID string, signal syscall.Signal) (err error) {
	conn, err := grpc.Dial(c.Endpoint.URL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := containers.NewContainersClient(conn)

	_, err = client.Signal(c.ctx, &containers.SignalRequest{
		Namespace:   c.Namespace,
		ContainerID: containerID,
		Signal:      int32(signal),
	})

	return err
}
