package controller

import (
	"testing"
	"time"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSyncStartsMultiContainerPod(t *testing.T) {
	in := make(chan []model.Pod)
	out := make(chan []model.Pod)
	clientMock := &FakeClient{t, []string{"default"}, map[string]map[string][]FakeContainer{}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			Metadata: model.Metadata{
				Name:      "hello-world",
				Namespace: "test",
			},
			Spec: model.PodSpec{
				Containers: []model.Container{
					model.Container{
						Name:  "hello-world-first",
						Image: "docker.io/eaapa/hello-world:latest",
					},
					model.Container{
						Name:  "hello-world-second",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := New(clientMock, 1*time.Second, in, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(2, 2, 0)
}

func TestSyncStopRemovedPodContainers(t *testing.T) {
	in := make(chan []model.Pod)
	out := make(chan []model.Pod)
	clientMock := &FakeClient{t, []string{"default", "cand"}, map[string]map[string][]FakeContainer{
		"cand": map[string][]FakeContainer{
			"my-pod": []FakeContainer{
				fakeRunningContainer("container-name", "docker.io/eaapa/hello-world:latest"),
			},
			"other-pod": []FakeContainer{
				fakeRunningContainer("will-be-removed", "docker.io/eaapa/hello-world:latest"),
			},
		},
	}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			Metadata: model.Metadata{
				Namespace: "cand",
				Name:      "my-pod",
			},
			Spec: model.PodSpec{
				Containers: []model.Container{
					model.Container{
						Name:  "container-name",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := New(clientMock, 1*time.Second, in, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 0, 1)
}

func TestSyncStartsMissingContainerTask(t *testing.T) {
	in := make(chan []model.Pod)
	out := make(chan []model.Pod)
	clientMock := &FakeClient{t, []string{"default", "cand"}, map[string]map[string][]FakeContainer{
		"cand": map[string][]FakeContainer{
			"my-pod": []FakeContainer{
				fakeCreatedContainer("container-name", "docker.io/eaapa/hello-world:latest"),
			},
		},
	}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			Metadata: model.Metadata{
				Namespace: "cand",
				Name:      "my-pod",
			},
			Spec: model.PodSpec{
				Containers: []model.Container{
					model.Container{
						Name:  "container-name",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := New(clientMock, 1*time.Second, in, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 1, 0)
}
