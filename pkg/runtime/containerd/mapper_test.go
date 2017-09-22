package containerd

import (
	"testing"

	"github.com/containerd/containerd"
	"github.com/stretchr/testify/assert"
)

func TestReconstructPods(t *testing.T) {
	containers := []containerd.Container{
		fakeRunningContainer("cand", "my-pod", "container1"),
		fakeRunningContainer("cand", "my-pod", "container2"),
		fakeCreatedContainer("cand", "my-other-pod", "hello-world-cont"),
	}

	result := MapToModelByPodNames(containers)

	assert.Len(t, result, 2, "Should construct two Pods from container information")
	assert.Len(t, result["my-pod"], 2, "Should have two containers for pod 'my-pod'")
	assert.Len(t, result["my-other-pod"], 1, "Should have two containers for pod 'my-pod'")

	assert.Equal(t, result["my-other-pod"][0].Name, "hello-world-cont", "Should have two containers for pod 'my-pod'")
}
