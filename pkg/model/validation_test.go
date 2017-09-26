package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationIsValid(t *testing.T) {
	err := Validate([]Pod{
		Pod{
			Metadata: Metadata{
				Name: "foo",
			},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name:  "foo-1",
						Image: "docker.io/library/foobar",
					},
				},
			},
		},
	})

	assert.NoError(t, err, "should be valid")
}

func TestValidationRequiresSpec(t *testing.T) {
	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{Name: "foo"},
		},
	}), "should return error if no spec defined")

	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{Name: "foo"},
			Spec:     PodSpec{},
		},
	}), "should return error if no containers defined")

	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{Name: "foo"},
			Spec: PodSpec{
				Containers: []Container{},
			},
		},
	}), "should return error if no any container defined")
}

func TestValidationRequiresContainerFields(t *testing.T) {
	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{
				Name: "foo",
			},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Image: "foobar",
					},
				},
			},
		},
	}), "should return error if container don't have name field")

	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{
				Name: "foo",
			},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name: "foo",
					},
				},
			},
		},
	}), "should return error if container don't have image field")

	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{
				Name: "foo",
			},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name:  "foo",
						Image: "/foo",
					},
				},
			},
		},
	}), "should return error if container image reference is invalid")
}

func TestValidationRequiresMetadata(t *testing.T) {
	assert.Error(t, Validate([]Pod{
		Pod{
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name:  "foo",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}), "should return error if no metadata")
}

func TestValidationNameMetadata(t *testing.T) {

	assert.Error(t, Validate([]Pod{
		Pod{
			Metadata: Metadata{},
			Spec: PodSpec{
				Containers: []Container{
					Container{
						Name:  "foo",
						Image: "docker.io/eaapa/hello-world:latest",
					},
				},
			},
		},
	}), "should return error if no 'name' metadata")

	assert.Error(t, Validate([]Pod{
		Pod{
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
		},
	}), "should return error if not alphanumeric name")
}

func TestValidationNamespaceMetadata(t *testing.T) {
	assert.Error(t, Validate([]Pod{
		Pod{
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
		},
	}), "should return error if not alphanumeric namespace")
}

func TestImageReferenceValidation(t *testing.T) {
	assert.True(t, isValidImageReference("docker.io/library/hello-world:latest"), "should be valid full image reference")
	assert.True(t, isValidImageReference("docker.io/library/hello-world"), "should be valid image reference without tag")
	assert.False(t, isValidImageReference("/hello-world"), "should be invalid reference if no hostname")
}
