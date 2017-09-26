package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodGetName(t *testing.T) {
	pod := Pod{
		Metadata: Metadata{
			Name: "foobar",
		},
	}

	assert.Equal(t, pod.Metadata.Name, "foobar", "should return name from metadata")
}

func TestPodGetNamespace(t *testing.T) {
	pod := Pod{
		Metadata: Metadata{
			Namespace: "foobar",
		},
	}

	assert.Equal(t, pod.Metadata.Namespace, "foobar", "should return namespace from metadata")
}

func TestContainerBuildID(t *testing.T) {
	assert.Equal(t, 36, len(BuildContainerID()), "should build valid id")
}
