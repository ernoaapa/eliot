package manifest

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ernoaapa/layeryd/model"
	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, data []byte) string {
	dir, createErr := ioutil.TempDir("", "example")
	assert.NoError(t, createErr, "Unable to create temp file")
	filePath := fmt.Sprintf("%s/%s", dir, "test.yml")

	writeErr := ioutil.WriteFile(filePath, data, 0666)
	assert.NoError(t, writeErr, "Unable to write to temporary file")
	return filePath
}

func TestFileSource(t *testing.T) {
	filePath := createTempFile(t, []byte(`
- metadata:
    name: "foo"
    namespace: "my-namespace"
  spec:
    containers:
      - name: "foo-1"
        image: "docker.io/library/hello-world:latest"
      - name: "foo-2"
        image: "docker.io/library/hello-world:latest"
- metadata:
    name: "bar"
  spec:
    containers:
      - name: "bar"
        image: "docker.io/library/hello-world:latest"
`))

	source := NewFileManifestSource(filePath, 100*time.Millisecond)
	updates := source.GetUpdates(model.DeviceInfo{})

	select {
	case pods := <-updates:
		assert.Equal(t, 2, len(pods), "Should have one pod spec")
		assert.Equal(t, "foo", pods[0].GetName(), "Should unmarshal name")
		assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")

		assert.Equal(t, "my-namespace", pods[0].GetNamespace(), "Should set default namespace")
		assert.Equal(t, "layeryd", pods[1].GetNamespace(), "Should set default namespace")
	case <-time.After(200 * time.Millisecond):
		assert.FailNow(t, "Didn't receive update in two second")
	}
}
