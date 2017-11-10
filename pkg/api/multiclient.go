package api

import (
	"fmt"
	"syscall"

	"golang.org/x/net/context"

	pods "github.com/ernoaapa/elliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/elliot/pkg/config"
	"github.com/ernoaapa/elliot/pkg/progress"
)

// MultiDirectClient connects directly to multiples devices RPC API
type MultiDirectClient struct {
	Namespace string
	clients   []Client
	ctx       context.Context
}

// NewMultiDirectClient creates new client which connects to multiple RPC APIs
func NewMultiDirectClient(namespace string, endpoints []config.Endpoint) *MultiDirectClient {
	clients := []Client{}
	for _, endpoint := range endpoints {
		clients = append(clients, NewDirectClient(namespace, endpoint))
	}
	return &MultiDirectClient{
		namespace,
		clients,
		context.Background(),
	}
}

// GetPods calls each device and fetches all pods information
func (c *MultiDirectClient) GetPods() (result []*pods.Pod, err error) {
	// TODO: optimise, run in parallel
	for _, client := range c.clients {
		pods, err := client.GetPods()
		if err != nil {
			return result, err
		}
		result = append(result, pods...)
	}
	return result, nil
}

// GetPod will return error, because getting pod from multiple devices is not yet supported
func (c *MultiDirectClient) GetPod(podName string) (*pods.Pod, error) {
	return nil, fmt.Errorf("Get pod from multiple devices is not yet supported. You need to specify device")
}

// CreatePod will return error, because creating pod to multiple devices is not yet supported
func (c *MultiDirectClient) CreatePod(updates chan<- []*progress.ImageFetch, pod *pods.Pod, opts ...PodOpts) error {
	return fmt.Errorf("Create pod to multiple devices is not yet supported. You need to specify device")
}

// StartPod will return error, because starting pod in multiple devices is not yet supported
func (c *MultiDirectClient) StartPod(name string) (*pods.Pod, error) {
	return nil, fmt.Errorf("Start pod in multiple devices is not yet supported. You need to specify device")
}

// DeletePod will return error, because deleting pod in multiple devices is not yet supported
func (c *MultiDirectClient) DeletePod(pod *pods.Pod) (*pods.Pod, error) {
	return nil, fmt.Errorf("Delete pod from multiple devices is not yet supported. You need to specify device")
}

// Attach will return error, because attaching to multiple devices is not yet supported
func (c *MultiDirectClient) Attach(containerID string, attachIO AttachIO, hooks ...AttachHooks) (err error) {
	return fmt.Errorf("Attach to multiple devices is not yet supported. You need to specify device")
}

// Signal will return error, because signaling to multiple devices is not yet supported
func (c *MultiDirectClient) Signal(containerID string, signal syscall.Signal) (err error) {
	return fmt.Errorf("Attach to multiple devices is not yet supported. You need to specify device")
}
