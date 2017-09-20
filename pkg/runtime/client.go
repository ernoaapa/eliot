package runtime

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
)

// Client interface for underlying implementation
type Client interface {
	GetContainers(namespace string) (containers []containerd.Container, err error)
	CreateContainer(pod model.Pod, container model.Container) error
	StopContainer(container containerd.Container) error
	EnsureImagePulled(namespace, ref string) (image containerd.Image, err error)
	GetNamespaces() ([]string, error)
	GetContainerTaskStatus(containerID string) string
}
