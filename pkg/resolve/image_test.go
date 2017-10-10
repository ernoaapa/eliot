package resolve

import (
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestImageResolve(t *testing.T) {
	projectDir := getExampleDirectory("node")
	image, err := Image(projectDir)
	assert.NoError(t, err)

	assert.Equal(t, "docker.io/library/node:latest", image)
}

func getExampleDirectory(name string) string {
	dir, err := filepath.Abs(filepath.Join(".", "examples", name))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
