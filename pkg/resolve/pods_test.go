package resolve

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestPodsResolveFile(t *testing.T) {
	examplePodFile := []byte(`
metadata:
  name: "hello-world"
spec:
  containers:
    - name: "hello-world"
      image: "docker.io/library/busybox:latest"
      args:
        - /bin/sh
        - -c
        - "while true; echo 'Eliot Rocks!'; do sleep 1; done;"
`)
	tmpfile, err := ioutil.TempFile("", "pods-resolve-test")
	assert.NoError(t, err)
	if ioutil.WriteFile(tmpfile.Name(), examplePodFile, 0644); err != nil {
		assert.Fail(t, "Failed go generate temp file: %s", err)
	}
	defer os.Remove(tmpfile.Name())

	result, err := Pods([]string{tmpfile.Name()})
	assert.NoError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, "hello-world", result[0].Metadata.Name)
	assert.Equal(t, "docker.io/library/busybox:latest", result[0].Spec.Containers[0].Image)
}
