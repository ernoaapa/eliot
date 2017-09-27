package runtime

import (
	"io"

	"github.com/ernoaapa/can/pkg/model"
)

// Client is interface for underlying container implementation
type Client interface {
	GetAllContainers(namespace string) (containersByPods map[string][]model.Container, err error)
	GetContainers(namespace, podName string) ([]model.Container, error)
	CreateContainer(pod model.Pod, container model.Container) error
	StartContainer(namespace, name string) error
	StopContainer(namespace, name string) error
	GetNamespaces() ([]string, error)
	IsContainerRunning(namespace, name string) (bool, error)
	GetContainerTaskStatus(namespace, name string) string
	Attach(namespace, podName string, attach AttachIO) error
}

// AttachIO provides way to attach stdin,stdout and stderr to container
type AttachIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NoIO() AttachIO {
	return AttachIO{}
}

// Empty returns true if the IO is missing one of io's
func (attach AttachIO) Empty() bool {
	return attach.Stdin == nil || attach.Stdout == nil || attach.Stderr == nil
}
