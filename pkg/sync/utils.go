package sync

import (
	"fmt"
	"log"
	"strings"
)

func MustParseAll(strings []string) (result []Sync) {
	for _, str := range strings {
		sync, err := Parse(str)
		if err != nil {
			log.Fatalf("Invalid formated sync [%s], must be in format: '<source>:<destination>', for example '~/local/dir:/data'", str)
		}
		result = append(result, sync)
	}
	return result
}

// Parse parses a sync string in the form "~/local/dir:/data"
func Parse(str string) (Sync, error) {
	parts := strings.Split(str, ":")

	if len(parts) == 2 {
		return Sync{
			Source:      parts[0],
			Destination: parts[1],
		}, nil
	}

	return Sync{}, fmt.Errorf("Invalid formated sync [%s], must be in format: '<source>:<destination>', for example '~/local/dir:/data'", str)
}
