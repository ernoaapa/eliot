package manifest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestUrlSource(t *testing.T) {
	updates := make(chan []model.Pod)
	resolver := device.NewResolver(map[string]string{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Should make PUT request")

		w.Header().Set(contentTypeHeader, yamlContentType)
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

	source := NewURLManifestSource(ts.URL, 100*time.Millisecond, resolver, updates)
	go source.Start()
	defer source.Stop()

	select {
	case pods := <-updates:
		assert.Equal(t, 2, len(pods), "Should have one pod spec")
		assert.Equal(t, "foo", pods[0].GetName(), "Should unmarshal name")
		assert.Equal(t, 2, len(pods[0].Spec.Containers), "Should have one container spec")

		assert.Equal(t, "my-namespace", pods[0].GetNamespace(), "Should set default namespace")
		assert.Equal(t, "cand", pods[1].GetNamespace(), "Should set default namespace")
	case <-time.After(200 * time.Millisecond):
		assert.FailNow(t, "Didn't receive update in two second")
	}
}

func TestUrlSourceHandlesUnauthorized(t *testing.T) {
	updates := make(chan []model.Pod)
	resolver := device.NewResolver(map[string]string{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set(contentTypeHeader, jsonContentType)
		w.Write([]byte("Whoops not authorized!"))

	}))
	defer ts.Close()

	source := NewURLManifestSource(ts.URL, 100*time.Millisecond, resolver, updates)
	go source.Start()
	defer source.Stop()

	select {
	case <-updates:
		assert.FailNow(t, "Where did you got that info?!")
	case <-time.After(200 * time.Millisecond):
		// Ok, didn't fail...
	}
}
