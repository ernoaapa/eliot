package controller

import (
	"context"
	"testing"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type FakeClient struct {
	t           *testing.T
	namespaces  []string
	containers  []containerd.Container
	createCount int
	stopCount   int
}

func (c *FakeClient) GetContainers(namespace string) (containers []containerd.Container, err error) {
	for _, container := range c.containers {
		if container.Info().Labels[runtime.GetLabelKeyFor(runtime.PodNamespaceSuffix)] == namespace {
			containers = append(containers, container)
		}
	}
	return containers, nil
}

func (c *FakeClient) CreateContainer(pod model.Pod, container model.Container) error {
	c.createCount++
	return nil
}

func (c *FakeClient) StopContainer(container containerd.Container) error {
	c.stopCount++
	return nil
}

func (c *FakeClient) EnsureImagePulled(namespace, ref string) (image containerd.Image, err error) {
	return image, nil
}

func (c *FakeClient) GetNamespaces() ([]string, error) {
	return c.namespaces, nil
}

func (c *FakeClient) GetContainerTaskStatus(containerID string) string {
	return ""
}

func (c *FakeClient) verifyExpectations(createCount, stopCount int) {
	assert.Equal(c.t, createCount, c.createCount, "Container create count should match")
	assert.Equal(c.t, stopCount, c.stopCount, "Container stop count should match")
}

type FakeContainer struct {
	id     string
	labels map[string]string
}

func (c *FakeContainer) ID() string {
	return c.id
}

func (c *FakeContainer) Info() containers.Container {
	return containers.Container{
		Labels: c.labels,
	}
}

func (c *FakeContainer) Delete(context.Context, ...containerd.DeleteOpts) error {
	return nil
}

func (c *FakeContainer) NewTask(context.Context, containerd.IOCreation, ...containerd.NewTaskOpts) (task containerd.Task, err error) {
	return task, err
}

func (c *FakeContainer) Spec() (*specs.Spec, error) {
	return nil, nil
}

func (c *FakeContainer) Task(context.Context, containerd.IOAttach) (task containerd.Task, err error) {
	return task, nil
}

func (c *FakeContainer) Image(context.Context) (image containerd.Image, err error) {
	return image, nil
}

func (c *FakeContainer) Labels(context.Context) (labels map[string]string, err error) {
	return labels, nil
}

func (c *FakeContainer) SetLabels(context.Context, map[string]string) (labels map[string]string, err error) {
	return labels, nil
}
func fakeRunningContainer(namespace, podName, containerName string) containerd.Container {
	uid := uuid.NewV4().String()

	labels := map[string]string{}
	labels[runtime.GetLabelKeyFor(runtime.PodUIDSuffix)] = uid
	labels[runtime.GetLabelKeyFor(runtime.PodNameSuffix)] = podName
	labels[runtime.GetLabelKeyFor(runtime.PodNamespaceSuffix)] = namespace
	labels[runtime.GetLabelKeyFor(runtime.ContainerNameSuffix)] = containerName

	return &FakeContainer{
		id:     uid,
		labels: labels,
	}
}
