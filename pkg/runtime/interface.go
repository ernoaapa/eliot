package runtime

import (
	"io"

	"github.com/ernoaapa/can/pkg/model"
)

// Client is interface for underlying container implementation
type Client interface {
	GetContainers(namespace string) (containersByPods map[string][]model.Container, err error)
	CreateContainer(pod model.Pod, container model.Container) error
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	GetNamespaces() ([]string, error)
	IsContainerRunning(containerID string) (bool, error)
	GetContainerTaskStatus(containerID string) string
	GetLogs(namespace, podName string, stdin io.Reader, stdout, stderr io.Writer) error
}
