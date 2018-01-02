package resolve

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ernoaapa/eliot/pkg/fs"
	log "github.com/sirupsen/logrus"
)

// Image try to resolve what container image should be used to run project in the directory
func Image(arch string, projectDir string) (projectType, image string, err error) {
	if isNodeProject(projectDir) {
		switch arch {
		case "amd64":
			return "nodejs", "docker.io/library/node:latest", nil
		case "arm64":
			return "nodejs", "docker.io/arm64v8/node:latest", nil
		default:
			return "", "", fmt.Errorf("Unsupported NodeJS project in architecture [%s]", arch)
		}
	} else if isGolangProject(projectDir) {
		switch arch {
		case "amd64":
			return "golang", "docker.io/library/golang:latest", nil
		case "arm64":
			return "golang", "docker.io/arm64v8/golang:latest", nil
		default:
			return "", "", fmt.Errorf("Unsupported Golang project in architecture [%s]", arch)
		}
	} else if isPythonProject(projectDir) {
		switch arch {
		case "amd64":
			return "python", "docker.io/library/python:latest", nil
		case "arm64":
			return "python", "docker.io/arm64v8/python:latest", nil
		default:
			return "", "", fmt.Errorf("Unsupported Golang project in architecture [%s]", arch)
		}
	}

	return "", "", fmt.Errorf("Unable to resolve container image for project in directory [%s]", projectDir)
}

func isNodeProject(projectDir string) bool {
	nodePackageFile := filepath.Join(projectDir, "package.json")
	log.Debugf("Checking does [%s] file exist, if does use Node container image", nodePackageFile)
	if fs.FileExist(nodePackageFile) {
		return true
	}
	return false
}

var golangDirs = []string{".", "pkg", "cmd"}

func isGolangProject(projectDir string) bool {
	for _, goDir := range golangDirs {
		if containsFiles(filepath.Join(projectDir, goDir), ".go") {
			return true
		}
	}
	return false
}

func isPythonProject(projectDir string) bool {
	for _, goDir := range golangDirs {
		if containsFiles(filepath.Join(projectDir, goDir), ".py") {
			return true
		}
	}
	return false
}

func containsFiles(path, extension string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), extension) {
			return true
		}
	}
	return false
}
