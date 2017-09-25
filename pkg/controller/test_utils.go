package controller

import (
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// FakeClient is runtime.Client implementation for tests to remove dependency to actual containerd
type FakeClient struct {
	t            *testing.T
	namespaces   []string
	containers   map[string]map[string][]FakeContainer
	createdCount int
	startedCount int
	stoppedCount int
}

// GetContainers fake impl.
func (c *FakeClient) GetContainers(namespace string) (map[string][]model.Container, error) {
	for podNamespace, podContainers := range c.containers {
		if podNamespace == namespace {
			result := map[string][]model.Container{}
			for podName, containers := range podContainers {
				result[podName] = []model.Container{}
				for _, fakeContainer := range containers {
					result[podName] = append(result[podName], fakeToModel(fakeContainer))
				}
			}
			return result, nil
		}
	}
	return make(map[string][]model.Container), nil
}

// CreateContainer fake impl.
func (c *FakeClient) CreateContainer(pod model.Pod, container model.Container) error {
	c.createdCount++
	return nil
}

// StartContainer fake impl.
func (c *FakeClient) StartContainer(containerID string) error {
	c.startedCount++
	return nil
}

// StopContainer fake impl.
func (c *FakeClient) StopContainer(containerID string) error {
	c.stoppedCount++
	return nil
}

// GetNamespaces fake impl.
func (c *FakeClient) GetNamespaces() ([]string, error) {
	return c.namespaces, nil
}

// IsContainerRunning fake impl.
func (c *FakeClient) IsContainerRunning(containerID string) (bool, error) {
	log.Debugf("container runnin %s", containerID)
	for _, podContainers := range c.containers {
		for _, containers := range podContainers {
			for _, fakeContainer := range containers {
				if fakeContainer.ID == containerID {
					log.Debugf("is running, %s, %s", containerID, fakeContainer.isRunning)
					return fakeContainer.isRunning, nil
				}
			}
		}
	}
	return false, nil
}

// GetContainerTaskStatus fake impl.
func (c *FakeClient) GetContainerTaskStatus(containerID string) string {
	return "UNKNOWN"
}

func (c *FakeClient) verifyExpectations(createdCount, startedCount, stoppedCount int) {
	assert.Equal(c.t, createdCount, c.createdCount, "Container create count should match")
	assert.Equal(c.t, startedCount, c.startedCount, "Container start count should match")
	assert.Equal(c.t, stoppedCount, c.stoppedCount, "Container stop count should match")
}

// FakeContainer is model.Container with some test related information, e.g. is it running
type FakeContainer struct {
	ID        string
	Name      string
	Image     string
	isRunning bool
}

func fakeRunningContainer(containerName, image string) FakeContainer {
	return newFakeContainer(containerName, image, true)
}

func fakeCreatedContainer(containerName, image string) FakeContainer {
	return newFakeContainer(containerName, image, false)
}

func newFakeContainer(containerName, image string, isRunning bool) FakeContainer {
	return FakeContainer{
		ID:        containerName,
		Name:      containerName,
		Image:     image,
		isRunning: isRunning,
	}
}

func fakeToModels(fakes []FakeContainer) (result []model.Container) {
	for _, fake := range fakes {
		result = append(result, fakeToModel(fake))
	}
	return result
}

func fakeToModel(fake FakeContainer) model.Container {
	return model.Container{
		ID:    fake.ID,
		Name:  fake.Name,
		Image: fake.Image,
	}
}
