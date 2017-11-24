package runtime

import (
	"io"
	"syscall"

	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/progress"
)

// Client is interface for underlying container implementation
type Client interface {
	GetPods(namespace string) ([]model.Pod, error)
	GetPod(namespace, podName string) (model.Pod, error)
	PullImage(namespace, ref string, status *progress.ImageFetch) error
	CreateContainer(pod model.Pod, container model.Container) (model.ContainerStatus, error)
	StartContainer(namespace, id string, io IOSet) (model.ContainerStatus, error)
	StopContainer(namespace, id string) (model.ContainerStatus, error)
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
