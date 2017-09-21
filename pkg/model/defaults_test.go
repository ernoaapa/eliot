package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	pods := []Pod{
		Pod{
			Metadata: Metadata{
				"name": "foobar",
			},
			Spec: PodSpec{
				Containers: []Container{},
			},
		},
		Pod{
			Metadata: Metadata{
				"name":      "foobar",
				"namespace": "my-namespace",
			},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name:  "foo",
						Image: "docker.io/library/hello-world:latest",
					},
				},
			},
		},
	}

	result := Defaults(pods)

	assert.Equal(t, "cand", result[0].GetNamespace(), "should set default namespace")
	assert.Equal(t, "my-namespace", result[1].GetNamespace(), "should not change namespace")

	assert.NotEmpty(t, result[1].Spec.Containers[0].ID, "should build container ID")
}
