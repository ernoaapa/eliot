package runtime

import (
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGetContainerLabels(t *testing.T) {
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
	result := getContainerLabels(pod, container)

	assert.Equal(t, "some-long-uid", result[getLabelKeyFor("pod.uid")])
	assert.Equal(t, "my-pod", result[getLabelKeyFor("pod.name")])
	assert.Equal(t, "my-namespace", result[getLabelKeyFor("pod.namespace")])
	assert.Equal(t, "my-container", result[getLabelKeyFor("container.name")])
}
