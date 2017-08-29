package manifest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUrlSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
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
`)
	}))
	defer ts.Close()

	source := NewURLManifestSource(ts.URL, 100*time.Millisecond)
	updates := source.GetUpdates()

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
