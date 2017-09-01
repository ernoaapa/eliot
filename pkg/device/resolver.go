package device

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Resolver provides information about the device
type Resolver struct {
	labels map[string]string
}

// NewResolver creates new resolver with static device labels
func NewResolver(labels map[string]string) *Resolver {
	return &Resolver{
		labels,
	}
}

func resolveFirst(name string, resolvers ...func() string) string {
	for _, resolver := range resolvers {
		result := resolver()
		if result != "" {
			return result
		}
	}

	log.Fatalf("Failed to resolve %s no default provided!", name)
	return ""
}

func fromEnv(name string) func() string {
	return func() string {
		return os.Getenv(name)
	}
}

func fromFiles(filePaths []string) func() string {
	return func() string {
		for _, file := range filePaths {
			info, err := ioutil.ReadFile(file)
			if err == nil {
				return strings.TrimSpace(string(info))
			}
		}
		return ""
	}
}
