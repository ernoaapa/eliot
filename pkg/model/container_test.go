package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidContainer(t *testing.T) {
	container := Container{
		Name:  "foo-1",
		Image: "docker.io/library/foobar",
		Args: []string{
			"/bin/bash",
			"ls",
		},
	}
	err := getValidator().Struct(container)
	assert.NoError(t, err, "should be valid")
}

func TestValidationContainerEnvVariables(t *testing.T) {
	assert.NoError(t, getValidator().Struct(Container{
		Name:  "foo-1",
		Image: "docker.io/library/foobar",
		Args: []string{
			"/bin/bash",
			"ls",
		},
		Env: []string{
			"FOO=bar",
		},
	}), "should be valid")
}
func TestValidationRequiresContainerFields(t *testing.T) {
	assert.Error(t, getValidator().Struct(Container{
		Image: "foobar",
	}), "should return error if container don't have name field")

	assert.Error(t, getValidator().Struct(Container{
		Name: "foo",
	}), "should return error if container don't have image field")

	assert.Error(t, getValidator().Struct(Container{
		Name:  "foo",
		Image: "/foo",
	}), "should return error if container image reference is invalid")
}
func TestContainerDevices(t *testing.T) {
	container := Container{
		Name:  "foo",
		Image: "docker.io/eaapa/hello-world:latest",
		Devices: []Device{
			Device{
				DeviceType : "c",
				MajorId: 555,
				MinorId: 0,
			},
		},
	}
	assert.Equal(t, len(container.Devices), 1, "should return 1")
}