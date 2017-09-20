package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodGetName(t *testing.T) {
	pod := Pod{
		Metadata: Metadata{
			"name": "foobar",
		},
	}

	assert.Equal(t, pod.GetName(), "foobar", "should return name from metadata")
}

func TestPodGetNamespace(t *testing.T) {
	pod := Pod{
		Metadata: Metadata{
			"namespace": "foobar",
		},
	}

	assert.Equal(t, pod.GetNamespace(), "foobar", "should return namespace from metadata")
}

func TestContainerBuildID(t *testing.T) {
	assert.Equal(t, 36, len(BuildContainerID()), "should build valid id")
}
