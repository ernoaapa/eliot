package utils

import (
	"fmt"
	"log"
	"strings"
)

var (
	defaultRegistry = "docker.io"
	defaultUsername = "library"
	defaultTag      = "latest"
)

// ExpandToFQIN converts partial image name to "Fully Qualified Image Name"
// E.g. eaapa/hello-world -> docker.io/eaapa/hello-world:latest
func ExpandToFQIN(source string) string {
	if source == "" {
		log.Fatal("Trying to expand empty image ref to FQIN (Fully Qualified Image Name)")
		return ""
	}
	registry := defaultRegistry
	username := defaultUsername
	tag := defaultTag
	image := source

	parts := strings.SplitN(source, "/", 3)
	if len(parts) == 3 {
		registry = parts[0]
		username = parts[1]
		image = parts[2]
	} else if len(parts) == 2 {
		username = parts[0]
		image = parts[1]
	}

	imageParts := strings.SplitN(image, ":", 2)
	if len(imageParts) == 2 {
		image = imageParts[0]
		tag = imageParts[1]
	}

	return fmt.Sprintf("%s/%s/%s:%s", registry, username, image, tag)
}

// GetFQINImage returns image part from FQIN
// E.g. docker.io/eaapa/hello-world:latest -> hello-world
func GetFQINImage(fqin string) string {
	parts := strings.SplitN(fqin, "/", 3)
	image := parts[2]
	imageParts := strings.SplitN(image, ":", 2)
	return imageParts[0]
}

// GetFQINUsername returns username part from FQIN
// E.g. docker.io/eaapa/hello-world:latest -> eaapa
func GetFQINUsername(fqin string) string {
	parts := strings.SplitN(fqin, "/", 3)
	return parts[1]
}
