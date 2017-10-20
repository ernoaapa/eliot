package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	file, tempErr := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, tempErr, "Failed to create temp file for test")
	defer os.Remove(file.Name())
	writeErr := ioutil.WriteFile(file.Name(), []byte(`
endpoints:
  - name: local-dev
    url: localhost:5000
namespace: foobar
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config, err := GetConfig(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Endpoints), "Should have one endpoint")
	assert.Equal(t, "foobar", config.Namespace, "Should return current context namespace")
}

func TestSetValue(t *testing.T) {
	config := Config{}

	config.Set("Namespace", "foobar")

	assert.Equal(t, "foobar", config.Namespace)
}
