package controller

import (
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSyncStartsMultiContainerPod(t *testing.T) {
	out := make(chan []model.Pod)
	clientMock := &FakeClient{t, []string{"default"}, map[string]map[string][]FakeContainer{}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			UID: "1",
			Metadata: model.Metadata{
				"name":      "hello-world",
				"namespace": "test",
			},
			Spec: model.PodSpec{
				Containers: []model.Container{
					model.Container{
						ID:    "hello-world-first",
						Name:  "hello-world-first",
						Image: "docker.io/eaapa/hello-world:latest",
					},
					model.Container{
						ID:    "hello-world-second",
						Name:  "hello-world-second",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := New(clientMock, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(2, 2, 0)
}

func TestSyncStopRemovedPodContainers(t *testing.T) {
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
			UID: "1",
			Metadata: model.Metadata{
				"namespace": "cand",
				"name":      "my-pod",
			},
			Spec: model.PodSpec{
				Containers: []model.Container{
					model.Container{
						ID:    "container-name",
						Name:  "container-name",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := New(clientMock, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 0, 1)
}

func TestSyncStartsMissingContainerTask(t *testing.T) {
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
			UID: "1",
			Metadata: model.Metadata{
				"namespace": "cand",
				"name":      "my-pod",
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

	err := New(clientMock, out).Sync(pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 1, 0)
}
