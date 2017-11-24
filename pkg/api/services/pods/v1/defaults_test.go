package pods

import (
	"testing"

	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	pods := []*Pod{
		&Pod{
			Metadata: &core.ResourceMetadata{
				Name: "foobar",
			},
			Spec: &PodSpec{
				Containers: []*containers.Container{},
			},
		},
		&Pod{
			Metadata: &core.ResourceMetadata{
				Name:      "foobar",
				Namespace: "my-namespace",
			},
			Spec: &PodSpec{
				Containers: []*containers.Container{
					&containers.Container{
						Name:  "foo",
						Image: "docker.io/library/hello-world:latest",
					},
				},
			},
		},
	}

	result := Defaults(pods)

	assert.Equal(t, "eliot", result[0].Metadata.Namespace, "should set default namespace")
	assert.Equal(t, "my-namespace", result[1].Metadata.Namespace, "should not change namespace")
	assert.NotEmpty(t, "foo", result[1].Spec.Containers[0].Name, "should not change name")
}
