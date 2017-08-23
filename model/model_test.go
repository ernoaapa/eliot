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

func TestContainerBuildID(t *testing.T) {
	container := Container{
		Name: "foo",
	}

	assert.Equal(t, container.BuildID("podname"), "podname-foo", "should build valid id")
}
