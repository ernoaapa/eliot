package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectConfig(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
name: foobar
image: someproject/foobar:latest
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config := ReadProjectConfig(file.Name())

	assert.Equal(t, "foobar", config.Name)
	assert.Equal(t, "someproject/foobar:latest", config.Image)
	assert.Equal(t, "docker.io/ernoaapa/rsync:1940a6c", config.Sync.Image, "sync.image should have default value")
}

func TestGetProjectSyncConfig(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
name: foobar
image: someproject/foobar:latest
sync:
  image: "myproject/custom:v1"
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config := ReadProjectConfig(file.Name())

	assert.Equal(t, "foobar", config.Name)
	assert.Equal(t, "someproject/foobar:latest", config.Image)
	assert.Equal(t, "myproject/custom:v1", config.Sync.Image)
}
