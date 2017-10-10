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
						Args: []string{
							"/bin/bash",
							"ls",
						},
					},
				},
			},
		},
	})

	assert.NoError(t, err, "should be valid")
}

func TestImageReferenceValidation(t *testing.T) {
	assert.True(t, isValidImageReference("docker.io/library/hello-world:latest"), "should be valid full image reference")
	assert.True(t, isValidImageReference("docker.io/library/hello-world"), "should be valid image reference without tag")
	assert.False(t, isValidImageReference("/hello-world"), "should be invalid reference if no hostname")
}

func TestEnvKeyValuePairs(t *testing.T) {
	assert.True(t, IsValidEnvKeyValuePair("FOO=bar"), "Should be valid env key/value pair")
	assert.True(t, IsValidEnvKeyValuePair("VERSION=12345"), "Should be valid env key/value pair")
	assert.True(t, IsValidEnvKeyValuePair("DEBUG=true"), "Should be valid env key/value pair")
	assert.True(t, IsValidEnvKeyValuePair("BAZ"), "Should be valid env key/value pair")

	assert.False(t, IsValidEnvKeyValuePair("%&%,foo"), "Should be invalid env key/value pair")
}
