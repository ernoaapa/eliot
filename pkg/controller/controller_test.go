package controller

import (
	"testing"

	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestSyncStartsMultiContainerPod(t *testing.T) {
	clientMock := &FakeClient{t, []string{"default"}, []containerd.Container{}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			UID: "1",
			Metadata: model.Metadata{
				"name":      "hello-world",
				"namespace": "test",
			},
			Spec: model.Spec{
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

	err := Sync(clientMock, pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(2, 2, 0)
}

func TestSyncStopRemovedContainers(t *testing.T) {
	clientMock := &FakeClient{t, []string{"default", "cand"}, []containerd.Container{
		fakeRunningContainer("cand", "my-pod", "container-name"),
		fakeRunningContainer("cand", "other-pod", "will-be-removed"),
	}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			UID: "1",
			Metadata: model.Metadata{
				"namespace": "cand",
				"name":      "my-pod",
			},
			Spec: model.Spec{
				Containers: []model.Container{
					model.Container{
						Name:  "container-name",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := Sync(clientMock, pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 0, 1)
}

func TestSyncStartsMissingContainerTask(t *testing.T) {
	clientMock := &FakeClient{t, []string{"default", "cand"}, []containerd.Container{
		fakeCreatedContainer("cand", "my-pod", "container-name"),
	}, 0, 0, 0}
	pods := []model.Pod{
		model.Pod{
			UID: "1",
			Metadata: model.Metadata{
				"namespace": "cand",
				"name":      "my-pod",
			},
			Spec: model.Spec{
				Containers: []model.Container{
					model.Container{
						Name:  "container-name",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}

	err := Sync(clientMock, pods)

	assert.NoError(t, err, "Sync should not return error")

	clientMock.verifyExpectations(0, 1, 0)
}
