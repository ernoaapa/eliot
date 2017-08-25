package device

import (
	"io/ioutil"
	"log"
	"strings"
)

func getInfoFromFiles(filePaths []string, defaultFn func([]string) string) string {
	for _, file := range filePaths {
		info, err := ioutil.ReadFile(file)
		if err == nil {
			return strings.TrimSpace(string(info))
		}
	}

	return defaultFn(filePaths)
}

func failIfCannotResolve(fieldName string) func([]string) string {
	return func(filePaths []string) string {
		log.Fatalf("Unable to resolve %s from following locations: %s", fieldName, strings.Join(filePaths, ","))
		return ""
	}
}
