package config

import (
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ernoaapa/can/pkg/converter"
	"github.com/ernoaapa/can/pkg/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config is struct for configuration for CLI client
type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
	Namespace string     `yaml:"namespace"`
}

// Endpoint represents Can server
type Endpoint struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// GetHost return just hostname/ip of endpoint URL
func (e Endpoint) GetHost() string {
	parts := strings.SplitN(e.URL, ":", 2)
	return parts[0]
}

// GetConfig reads current config from user home directory
func GetConfig(path string) (*Config, error) {
	config := &Config{
		Namespace: "cand",
	}
	if !fs.FileExist(path) {
		return config, nil
	}

	data, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return nil, errors.Wrapf(readErr, "Failed to read configuration file: %s", path)
	}

	unmarshalErr := yaml.Unmarshal(data, &config)
	if unmarshalErr != nil {
		return nil, errors.Wrapf(unmarshalErr, "Unable to parse Yaml config: %s", path)
	}

	return config, nil
}

// WriteConfig writes configuration to given path
func WriteConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "Unable to marshal config to yaml format")
	}
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return errors.Wrapf(err, "Failed to write config to path [%s]", path)
	}

	return nil
}

// Set mutates the Config by updating single field with value
func (c *Config) Set(field, value string) {
	v := reflect.ValueOf(c).Elem().FieldByName(converter.KebabCaseToCamelCase(field))
	if v.IsValid() {
		v.SetString(value)
	}
}
