package containerd

import (
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewContainerLabels(t *testing.T) {
	pod := model.Pod{
		UID: "some-long-uid",
		Metadata: model.Metadata{
			"name":      "my-pod",
			"namespace": "my-namespace",
		},
	}
	container := model.Container{
		Name: "my-container",
	}
	result := NewContainerLabels(pod, container)

	assert.Equal(t, "some-long-uid", result["io.can.pod.uid"])
	assert.Equal(t, "my-pod", result["io.can.pod.name"])
	assert.Equal(t, "my-namespace", result["io.can.pod.namespace"])
	assert.Equal(t, "my-container", result["io.can.container.name"])
}
