package device

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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
