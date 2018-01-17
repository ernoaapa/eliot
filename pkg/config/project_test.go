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
	assert.NotEmpty(t, config.SyncContainer.Image, "syncContainer.image should have default value")
}

func TestGetProjectSyncConfig(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
name: foobar
image: someproject/foobar:latest
syncContainer:
  image: "myproject/custom:v1"
syncs:
    - out:/app
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config := ReadProjectConfig(file.Name())

	assert.Equal(t, "foobar", config.Name)
	assert.Equal(t, "someproject/foobar:latest", config.Image)
	assert.Equal(t, "myproject/custom:v1", config.SyncContainer.Image)
	assert.Equal(t, []string{"out:/app"}, config.Syncs)
}

func TestGetProjectMountConfig(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
name: foobar
image: someproject/foobar:latest
mounts:
    - source=/dev,destination=/host-dev
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config := ReadProjectConfig(file.Name())

	assert.Equal(t, "foobar", config.Name)
	assert.Equal(t, "someproject/foobar:latest", config.Image)
	assert.Equal(t, []string{"source=/dev,destination=/host-dev"}, config.Mounts)
}

func TestGetProjectCommand(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
name: foobar
image: someproject/foobar:latest
command: ["foo", "bar"]
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config := ReadProjectConfig(file.Name())

	assert.Equal(t, "foobar", config.Name)
	assert.Equal(t, []string{"foo", "bar"}, config.Command)
}
