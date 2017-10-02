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

func TestValidationRequiresSpec(t *testing.T) {
	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{Name: "foo"},
	}), "should return error if no spec defined")

	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{Name: "foo"},
		Spec:     PodSpec{},
	}), "should return error if no containers defined")

	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{Name: "foo"},
		Spec: PodSpec{
			Containers: []Container{},
		},
	}), "should return error if no any container defined")
}

func TestValidationRequiresMetadata(t *testing.T) {
	assert.Error(t, getValidator().Struct(Pod{
		Spec: PodSpec{
			Containers: []Container{
				Container{
					Name:  "foo",
					Image: "docker.io/eaapa/hello-world:latest",
				},
			},
		},
	}), "should return error if no metadata")
}

func TestValidationNameMetadata(t *testing.T) {
	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{},
		Spec: PodSpec{
			Containers: []Container{
				Container{
					Name:  "foo",
					Image: "docker.io/eaapa/hello-world:latest",
				},
			},
		},
	}), "should return error if no 'name' metadata")

	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{
			Name: "#€%&/()=",
		},
		Spec: PodSpec{
			Containers: []Container{
				Container{
					Name:  "foo",
					Image: "docker.io/eaapa/hello-world:latest",
				},
			},
		},
	}), "should return error if not alphanumeric name")
}

func TestValidationNamespaceMetadata(t *testing.T) {
	assert.Error(t, getValidator().Struct(Pod{
		Metadata: Metadata{
			Name:      "foo",
			Namespace: "#€%&/()=",
		},
		Spec: PodSpec{
			Containers: []Container{
				Container{
					Name:  "foo",
					Image: "docker.io/eaapa/hello-world:latest",
				},
			},
		},
	}), "should return error if not alphanumeric namespace")
}
