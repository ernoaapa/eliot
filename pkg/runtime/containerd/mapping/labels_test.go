package mapping

import (
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewContainerLabels(t *testing.T) {
	pod := model.Pod{
		Metadata: model.Metadata{
			Name:      "my-pod",
			Namespace: "my-namespace",
		},
	}
	container := model.Container{
		Name: "my-container",
	}
	result := NewLabels(pod, container)

	assert.Equal(t, "my-pod", result["io.can.pod.name"])
}
