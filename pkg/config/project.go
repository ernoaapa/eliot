package config

import (
	"io/ioutil"
	"log"
	"reflect"

	"github.com/ernoaapa/can/pkg/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ProjectConfig represents configuration in project directory
type ProjectConfig struct {
	path    string
	configs map[string]interface{}
}

// ReadProjectConfig read ProjectConfig from given path
// If file doesn't exist, returns empty config so it's safe for reading even
// The file doesn't exist. In any other failure case, will fatal.
func ReadProjectConfig(path string) *ProjectConfig {
	configs := map[string]interface{}{}

	if !fs.FileExist(path) {
		return &ProjectConfig{path, configs}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "Failed to read project configuration at %s", path))
	}

	if err := yaml.Unmarshal(data, configs); err != nil {
		log.Fatalln(errors.Wrapf(err, "Error while reading YAML configuration at %s", path))
	}

	return &ProjectConfig{path, configs}
}

// String return string value from configuration or defaults to other value
// In case of invalid type fatal
func (c *ProjectConfig) String(key, defaultValue string) string {
	val, ok := c.configs[key]
	if !ok {
		return defaultValue
	}

	s, ok := val.(string)
	if !ok {
		log.Fatalf("Config at %s contains invalid value at %s. It should be string", c.path, key)
	}

	return s
}

// StringSlice return string slice value from configuration
// Will append extra fields to the result
// In case of invalid type fatal
func (c *ProjectConfig) StringSlice(key string, extra []string) []string {
	val, ok := c.configs[key]
	if !ok {
		return extra
	}

	s, ok := toStringSlice(val)
	if !ok {
		log.Fatalf("Config at %s contains invalid value at %s. It should be list of strings", c.path, key)
	}

	return append(s, extra...)
}

func toStringSlice(val interface{}) (result []string, ok bool) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(val)

		for i := 0; i < s.Len(); i++ {
			str, ok := s.Index(i).Interface().(string)
			if !ok {
				return result, false
			}
			result = append(result, str)
		}
	default:
		return result, false
	}
	return result, true
}
