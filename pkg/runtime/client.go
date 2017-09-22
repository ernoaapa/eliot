package runtime

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
)

// Client interface for underlying implementation
type Client interface {
	GetContainers(namespace string) (containers []containerd.Container, err error)
	CreateContainer(pod model.Pod, container model.Container) (containerd.Container, error)
	StartContainer(container containerd.Container) error
	StopContainer(container containerd.Container) error
	GetNamespaces() ([]string, error)
	IsContainerRunning(container model.Container) (bool, error)
	GetContainerTaskStatus(containerID string) string
}
