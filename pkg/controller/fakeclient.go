package controller

import (
	"syscall"
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/progress"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

// FakeClient is runtime.Client implementation for tests to remove dependency to actual containerd
type FakeClient struct {
	t            *testing.T
	namespaces   []string
	pods         map[string][]model.Pod
	createdCount int
	startedCount int
	stoppedCount int
}

// GetPods fake impl.
func (c *FakeClient) GetPods(namespace string) (result []model.Pod, err error) {
	return c.pods[namespace], nil
}

// GetPod fake impl.
func (c *FakeClient) GetPod(namespace, podName string) (result model.Pod, err error) {
	for podNamespace, pods := range c.pods {
		if podNamespace == namespace {
			for _, pod := range pods {
				if pod.Metadata.Name == podName {
					return pod, nil
				}
			}
		}
	}
	return result, err
}

// PullImage fake impl.
func (c *FakeClient) PullImage(namespace, ref string, progress *progress.ImageFetch) error {
	return nil
}

// CreateContainer fake impl.
func (c *FakeClient) CreateContainer(pod model.Pod, container model.Container) error {
	c.createdCount++
	return nil
}

// StartContainer fake impl.
func (c *FakeClient) StartContainer(namespace, containerID string, tty bool) (model.ContainerStatus, error) {
	c.startedCount++
	return model.ContainerStatus{
		ContainerID: containerID,
		State:       "running",
	}, nil
}

// StopContainer fake impl.
func (c *FakeClient) StopContainer(namespace, containerID string) (model.ContainerStatus, error) {
	c.stoppedCount++
	return model.ContainerStatus{
		ContainerID: containerID,
		State:       "stopped",
	}, nil
}

// GetNamespaces fake impl.
func (c *FakeClient) GetNamespaces() ([]string, error) {
	return c.namespaces, nil
}

// IsContainerRunning fake impl.
func (c *FakeClient) IsContainerRunning(namespace, name string) (bool, error) {
	for podNamespace, pods := range c.pods {
		for _, pod := range pods {
			if podNamespace == namespace {
				for _, containerStatus := range pod.Status.ContainerStatuses {
					if containerStatus.ContainerID == name {
						return containerStatus.State == "running", nil
					}
				}
			}
		}
	}
	return false, nil
}

// GetContainerTaskStatus fake impl.
func (c *FakeClient) GetContainerTaskStatus(namespace, name string) string {
	return "UNKNOWN"
}

// Attach fake impl.
func (c *FakeClient) Attach(namespace, podName string, io runtime.AttachIO) error {
	return nil
}

// Signal fake impl.
func (c *FakeClient) Signal(namespace, name string, signal syscall.Signal) error {
	return nil
}

func (c *FakeClient) verifyExpectations(createdCount, startedCount, stoppedCount int) {
	assert.Equal(c.t, createdCount, c.createdCount, "Container create count should match")
	assert.Equal(c.t, startedCount, c.startedCount, "Container start count should match")
	assert.Equal(c.t, stoppedCount, c.stoppedCount, "Container stop count should match")
}

type createOpts func(*model.Pod)

func newPod(namespace, name string, opts ...createOpts) model.Pod {
	pod := model.NewPod(name, namespace)
	for _, opt := range opts {
		opt(&pod)
	}
	return pod
}

func withRunningContainer(containerName, image string) createOpts {
	return func(pod *model.Pod) {
		pod.Spec.Containers = append(pod.Spec.Containers, model.Container{
			Name:  containerName,
			Image: image,
		})

		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, model.ContainerStatus{
			ContainerID: containerName,
			Image:       image,
			State:       "running",
		})
	}
}

func withCreatedContainer(containerName, image string) createOpts {
	return func(pod *model.Pod) {
		pod.Spec.Containers = append(pod.Spec.Containers, model.Container{
			Name:  containerName,
			Image: image,
		})

		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, model.ContainerStatus{
			ContainerID: containerName,
			Image:       image,
			State:       "created",
		})
	}
}
