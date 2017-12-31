package pods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalYaml(t *testing.T) {
	pods, err := UnmarshalYaml([]byte(`
metadata:
  name: "foo"
  namespace: "my-namespace"
spec:
  hostNetwork: true
  restartPolicy: always
  containers:
    - name: "foo-1"
      image: "docker.io/library/hello-world:latest"
    - name: "foo-2"
      image: "docker.io/library/hello-world:latest"
`))

	assert.NoError(t, err, "Unable unmarshal test yaml")

	assert.Equal(t, "foo", pods[0].Metadata.Name, "Should unmarshal name")
	assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")
	assert.Equal(t, "always", pods[0].Spec.RestartPolicy, "Should have host network")
	assert.True(t, pods[0].Spec.HostNetwork, "Should have host network")
}

func TestUnmarshalMultiDocumentYaml(t *testing.T) {
	pods, err := UnmarshalYaml([]byte(`
metadata:
  name: "foo"
  namespace: "my-namespace"
spec:
  containers:
    - name: "foo-1"
      image: "docker.io/library/hello-world:latest"
    - name: "foo-2"
      image: "docker.io/library/hello-world:latest"
---
metadata:
  name: "bar"
spec:
  containers:
    - name: "bar"
      image: "docker.io/library/hello-world:latest"
`))

	assert.NoError(t, err, "Unable unmarshal test yaml")
	assert.Equal(t, 2, len(pods), "Should have pod specs")
}

func TestUnmarshalListYaml(t *testing.T) {
	pods, err := UnmarshalListYaml([]byte(`
- metadata:
    name: "foo"
    namespace: "my-namespace"
  spec:
    hostNetwork: true
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

	assert.NoError(t, err, "Unable unmarshal test yaml")

	assert.Equal(t, 2, len(pods), "Should have pod specs")
	assert.Equal(t, "foo", pods[0].Metadata.Name, "Should unmarshal name")
	assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")
	assert.True(t, pods[0].Spec.HostNetwork, "Should have one container spec")
}

func TestUnmarshalListJSON(t *testing.T) {
	pods, err := UnmarshalListJSON([]byte(`
[
  {
    "metadata": {
      "name": "foo",
      "namespace": "my-namespace"
    },
    "spec": {
      "containers": [
        {
          "name": "foo-1",
          "image": "docker.io/library/hello-world:latest"
        },
        {
          "name": "foo-2",
          "image": "docker.io/library/hello-world:latest"
        }
      ]
    }
  },
  {
    "metadata": {
      "name": "bar"
    },
    "spec": {
      "containers": [
        {
          "name": "bar",
          "image": "docker.io/library/hello-world:latest"
        }
      ]
    }
  }
]
`))

	assert.NoError(t, err, "Unable unmarshal test yaml")

	assert.Equal(t, 2, len(pods), "Should have pod specs")
	assert.Equal(t, "foo", pods[0].Metadata.Name, "Should unmarshal name")
	assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")
}
