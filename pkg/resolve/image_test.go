package resolve

import (
	"fmt"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var targetArchitectures = []string{"amd64", "arm64"}

func TestImageResolveNode(t *testing.T) {
	projectDir := getExampleDirectory("node")

	testResolving(t, projectDir, "nodejs", map[string]string{
		"amd64": "docker.io/library/node:latest",
		"arm64": "docker.io/arm64v8/node:latest",
	})
}

func TestImageResolveGolang(t *testing.T) {
	projectDir := getExampleDirectory("golang")

	testResolving(t, projectDir, "golang", map[string]string{
		"amd64": "docker.io/library/golang:latest",
		"arm64": "docker.io/arm64v8/golang:latest",
	})
}

func TestImageResolvePython(t *testing.T) {
	projectDir := getExampleDirectory("python")

	testResolving(t, projectDir, "python", map[string]string{
		"amd64": "docker.io/library/python:latest",
		"arm64": "docker.io/arm64v8/python:latest",
	})
}

func getExampleDirectory(name string) string {
	dir, err := filepath.Abs(filepath.Join(".", "examples", name))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func testResolving(t *testing.T, projectDir, expectedProjectType string, images map[string]string) {
	for _, arch := range targetArchitectures {
		expectedImage, ok := images[arch]
		if !ok {
			assert.FailNow(t, fmt.Sprintf("Test case is missing expected image for projectType %s and architecture %s", expectedProjectType, arch))
		}

		projectType, image, err := Image(arch, projectDir)
		assert.NoError(t, err)
		assert.Equal(t, expectedProjectType, projectType)
		assert.Equal(t, expectedImage, image)
	}
}
