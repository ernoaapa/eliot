package runtime

import (
	"io"
	"syscall"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/progress"
)

// Client is interface for underlying container implementation
type Client interface {
	GetAllContainers(namespace string) (containersByPods map[string][]model.Container, err error)
	GetContainers(namespace, podName string) ([]model.Container, error)
	PullImage(namespace, ref string, status *progress.ImageFetch) error
	CreateContainer(pod model.Pod, container model.Container) error
	StartContainer(namespace, name string, tty bool) error
	StopContainer(namespace, name string) error
	GetNamespaces() ([]string, error)
	IsContainerRunning(namespace, name string) (bool, error)
	GetContainerTaskStatus(namespace, name string) string
	Attach(namespace, podName string, attach AttachIO) error
	Signal(namespace, name string, signal syscall.Signal) error
}

// AttachIO provides way to attach stdin,stdout and stderr to container
type AttachIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
