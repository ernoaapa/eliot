package build

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleLinuxkitConfig = []byte(`
kernel:
  image: linuxkit/kernel:4.9.70
`)

func TestResolveDefaultLinuxkitConfig(t *testing.T) {
	// Warning, requires github.com/ernoaapa/eliot-os access
	config, err := ResolveLinuxkitConfig("")
	assert.NoError(t, err)

	assert.True(t, len(config) > 0, "Should default to the default rpi3 config")
}

func TestResolveUrlLinuxkitConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(exampleLinuxkitConfig)
	}))
	defer ts.Close()
	config, err := ResolveLinuxkitConfig(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, exampleLinuxkitConfig, config, "Should fetch config from url")
}

func TestResolveFileLinuxkitConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	assert.NoError(t, err)
	if ioutil.WriteFile(tmpfile.Name(), exampleLinuxkitConfig, 0644); err != nil {
		assert.Fail(t, "Failed go generate temp file: %s", err)
	}

	defer os.Remove(tmpfile.Name())

	config, err := ResolveLinuxkitConfig(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, exampleLinuxkitConfig, config, "Should read config from file")
}

func TestBuildImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "fake image tar content")
	}))
	defer ts.Close()

	image, err := BuildImage(ts.URL, "rpi3", "img", exampleLinuxkitConfig)
	assert.NoError(t, err)

	tar, err := ioutil.ReadAll(image)
	assert.NoError(t, err)

	assert.Equal(t, tar, []byte("fake image tar content"))
}

func TestBuildImageReturnErrorMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprint(w, "This is the error")
	}))
	defer ts.Close()

	_, err := BuildImage(ts.URL, "rpi3", "tar", exampleLinuxkitConfig)
	assert.True(t, strings.Contains(err.Error(), "This is the error"))
}
