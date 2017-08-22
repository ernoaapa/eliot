package source

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/ernoaapa/layeryd/model"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	dir, createErr := ioutil.TempDir("", "example")
	assert.NoError(t, createErr, "Unable to create temp file")

	filePath := fmt.Sprintf("%s/%s", dir, "test.yml")
	log.Println(filePath)
	ioutil.WriteFile(filePath, []byte(`
  name: "foo"
  spec:
    containers:
      - name: "foo"
        image: "docker.io/library/hello-world:latest"
  `), 0666)

	source := NewFileSource(filePath)

	pod, err := source.GetState(model.NodeInfo{})
	assert.NoError(t, err, "Unable to read file")

	assert.Equal(t, "foo", pod.Name, "Should unmarshal name")
	// assert.Equal(t, 1, len(pod.Spec.Containers), "Should have one container spec")
}
