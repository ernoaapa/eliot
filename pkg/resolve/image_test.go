package resolve

import (
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestImageResolveNode(t *testing.T) {
	projectDir := getExampleDirectory("node")
	projectType, image, err := Image(projectDir)
	assert.NoError(t, err)

	assert.Equal(t, "nodejs", projectType)
	assert.Equal(t, "docker.io/library/node:latest", image)
}

func TestImageResolveGolang(t *testing.T) {
	projectDir := getExampleDirectory("golang")
	projectType, image, err := Image(projectDir)
	assert.NoError(t, err)

	assert.Equal(t, "golang", projectType)
	assert.Equal(t, "docker.io/library/golang:latest", image)
}

func getExampleDirectory(name string) string {
	dir, err := filepath.Abs(filepath.Join(".", "examples", name))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
