package resolve

import (
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var targetArchitectures = []string{"amd64", "arm64"}

func TestImageResolveNode(t *testing.T) {
	projectType, image, err := Image(getExampleDirectory("node"))
	assert.NoError(t, err)
	assert.Equal(t, "node", projectType)
	assert.Equal(t, "docker.io/library/node:latest", image)
}

func TestImageResolveGolang(t *testing.T) {
	projectType, image, err := Image(getExampleDirectory("golang"))
	assert.NoError(t, err)
	assert.Equal(t, "golang", projectType)
	assert.Equal(t, "docker.io/library/golang:latest", image)
}

func TestImageResolvePython(t *testing.T) {
	projectType, image, err := Image(getExampleDirectory("python"))
	assert.NoError(t, err)
	assert.Equal(t, "python", projectType)
	assert.Equal(t, "docker.io/library/python:latest", image)
}

func getExampleDirectory(name string) string {
	dir, err := filepath.Abs(filepath.Join(".", "examples", name))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
