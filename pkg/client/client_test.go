package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestClientGetDeployments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Should make GET request")
		assert.Equal(t, "Bearer the-token", r.Header.Get("Authorization"))

		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `
[
	{
    "metadata": {
      "name": "test-deployment"
    },
    "spec": {
      "selector": {
        "foo": "bar"
      },
      "template": {
        "metadata": {
          "name": "my-service"
        },
        "spec": {
          "containers": [{
            "name": "foo-1",
            "image": "docker.io/library/hello-world:latest"
          }]
        }
      }
    }
  }
]
`)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "the-token")

	deployments, err := client.GetDeployments()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(deployments), "Should have one deployment")
	assert.Equal(t, "test-deployment", deployments[0].GetName(), "Should have metadata.name")
	assert.Equal(t, 1, len(deployments[0].Spec.Template.Spec.Containers), "Should have one container spec in pod template")
}

func TestClientCreateDeployment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Should make GET request")
		assert.Equal(t, "application/json", r.Header.Get("content-type"))
		w.Header().Set("content-type", "application/json")
		io.Copy(w, r.Body)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "the-token")

	deployment := &model.Deployment{
		Metadata: model.Metadata{
			"name": "foobar",
		},
		Spec: model.DeploymentSpec{
			Selector: map[string]string{
				"foo": "bar",
			},
		},
	}
	created, err := client.CreateDeployment(deployment)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, deployment.GetName(), created.GetName(), "Should have same metadata.name")
}
