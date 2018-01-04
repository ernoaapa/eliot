package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandToFQIN(t *testing.T) {
	assert.Equal(t, "docker.io/eaapa/hello-world:latest", ExpandToFQIN("eaapa/hello-world"))
	assert.Equal(t, "otherhost.io/eaapa/hello-world:latest", ExpandToFQIN("otherhost.io/eaapa/hello-world"))
	assert.Equal(t, "docker.io/eaapa/hello-world:latest", ExpandToFQIN("eaapa/hello-world:latest"))
	assert.Equal(t, "docker.io/library/nginx:tag1", ExpandToFQIN("nginx:tag1"))
	assert.Equal(t, "docker.io/library/nginx:latest", ExpandToFQIN("nginx"))
}

func TestGetFQINImage(t *testing.T) {
	assert.Equal(t, "hello-world", GetFQINImage("docker.io/eaapa/hello-world"))
	assert.Equal(t, "hello-world", GetFQINImage("docker.io/eaapa/hello-world:latest"))
}

func TestGetFQINUsername(t *testing.T) {
	assert.Equal(t, "eaapa", GetFQINUsername("docker.io/eaapa/hello-world"))
	assert.Equal(t, "eaapa", GetFQINUsername("docker.io/eaapa/hello-world:latest"))
}
