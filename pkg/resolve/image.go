package resolve

import (
	"fmt"
	"path/filepath"

	"github.com/ernoaapa/can/pkg/fs"
	log "github.com/sirupsen/logrus"
)

// Image try to resolve what container image should be used to run project in the directory
func Image(projectDir string) (string, error) {
	nodePackageFile := filepath.Join(projectDir, "package.json")
	log.Debugf("Checking does [%s] file exist, if does use Node container image", nodePackageFile)
	if fs.FileExist(nodePackageFile) {
		return "docker.io/library/node:latest", nil
	}

	return "", fmt.Errorf("Unable to resolve container image for project in directory [%s]", projectDir)
}
