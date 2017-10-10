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
users:
  - name: erno
contexts:
  - name: local-dev/erno
    endpoint: local-dev
    user: erno
    namespace: foobar
current-context: local-dev/erno
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config, err := GetConfig(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Endpoints), "Should have one endpoint")
	assert.Equal(t, 1, len(config.Users), "Should have one user")
	assert.Equal(t, 1, len(config.Contexts), "Should have one context")
	assert.Equal(t, "local-dev/erno", config.CurrentContext, "Should parse active context")

	assert.Equal(t, "localhost:5000", config.GetCurrentEndpoint().URL, "Should return current user")
	assert.Equal(t, "foobar", config.GetCurrentContext().Namespace, "Should return current context namespace")
}

func TestSetValue(t *testing.T) {
	config := Config{}

	config.Set("CurrentContext", "foobar")

	assert.Equal(t, "foobar", config.CurrentContext)
}
