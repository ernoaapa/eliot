package source

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ernoaapa/layeryd/model"
	"github.com/stretchr/testify/assert"
)

func TestFileSource(t *testing.T) {
	dir, createErr := ioutil.TempDir("", "example")
	assert.NoError(t, createErr, "Unable to create temp file")
	filePath := fmt.Sprintf("%s/%s", dir, "test.yml")

	ioutil.WriteFile(filePath, []byte(`
- metadata:
    name: "foo"
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
`), 0666)

	source := NewFileSource(filePath, 100*time.Millisecond)
	updates := source.GetUpdates(model.DeviceInfo{})

	select {
	case pods := <-updates:
		assert.Equal(t, 2, len(pods), "Should have one pod spec")
		assert.Equal(t, "foo", pods[0].GetName(), "Should unmarshal name")
		assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")
	case <-time.After(200 * time.Millisecond):
		assert.FailNow(t, "Didn't receive update in two second")
	}
}
