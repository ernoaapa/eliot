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
    url: http://localhost:3000
users:
  - name: erno
    token: 123abcd
contexts:
  - name: local-dev/erno
    endpoint: local-dev
    user: erno
current-context: local-dev/erno
`), 0644)
	assert.NoError(t, writeErr, "Error while writing temp file")

	config, err := GetConfig(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Endpoints), "Should have one endpoint")
	assert.Equal(t, 1, len(config.Users), "Should have one endpoint")
	assert.Equal(t, 1, len(config.Contexts), "Should have one endpoint")
	assert.Equal(t, "local-dev/erno", config.CurrentContext, "Should parse active context")

	assert.Equal(t, "123abcd", config.GetCurrentUser().Token, "Should return current user")
	assert.Equal(t, "http://localhost:3000", config.GetCurrentEndpoint().URL, "Should return current user")
}
