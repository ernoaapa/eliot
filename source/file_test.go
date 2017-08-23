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
name: "foo"
spec:
  containers:
    - name: "foo"
      image: "docker.io/library/hello-world:latest"
`), 0666)

	source := NewFileSource(filePath, 100*time.Millisecond)
	updates := source.GetUpdates(model.DeviceInfo{})

	select {
	case pod := <-updates:
		assert.Equal(t, "foo", pod.Name, "Should unmarshal name")
		assert.Equal(t, 1, len(pod.Spec.Containers), "Should have one container spec")
	case <-time.After(200 * time.Millisecond):
		assert.FailNow(t, "Didn't receive update in two second")
	}
}
