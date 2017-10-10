package config

import (
	"io/ioutil"
	"reflect"

	"github.com/ernoaapa/can/pkg/converter"
	"github.com/ernoaapa/can/pkg/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config is struct for configuration for CLI client
type Config struct {
	Endpoints      []Endpoint `yaml:"endpoints"`
	Users          []User     `yaml:"users"`
	Contexts       []Context  `yaml:"contexts"`
	CurrentContext string     `yaml:"current-context"`
}

// Endpoint represents Can server
type Endpoint struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// User is authenticated user
type User struct {
	Name string `yaml:"name"`
}

// Context represents single user in single endpoint
type Context struct {
	Name      string `yaml:"name"`
	User      string `yaml:"user"`
	Endpoint  string `yaml:"endpoint"`
	Namespace string `yaml:"namespace"`
}

// DefaultConfig returns new Config instance with default values
func DefaultConfig() *Config {
	return &Config{
		Contexts: []Context{
			{Namespace: "cand"},
		},
	}
}

// GetConfig reads current config from user home directory
func GetConfig(path string) (*Config, error) {
	if !fs.FileExist(path) {
		return DefaultConfig(), nil
	}

	data, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return nil, errors.Wrapf(readErr, "Failed to read configuration file: %s", path)
	}

	config := &Config{}
	unmarshalErr := yaml.Unmarshal(data, config)
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

// GetCurrentUser return current user
func (c Config) GetCurrentUser() User {
	return c.GetUser(c.GetCurrentContext().User)
}

// GetCurrentContext return current context
func (c Config) GetCurrentContext() Context {
	return c.GetContext(c.CurrentContext)
}

// GetCurrentEndpoint return current context
func (c Config) GetCurrentEndpoint() Endpoint {
	return c.GetEndpoint(c.GetCurrentContext().Endpoint)
}

// Set mutates the Config by updating single field with value
func (c *Config) Set(field, value string) {
	v := reflect.ValueOf(c).Elem().FieldByName(converter.KebabCaseToCamelCase(field))
	if v.IsValid() {
		v.SetString(value)
	}
}

// GetContext return context by name
func (c Config) GetContext(name string) Context {
	for _, context := range c.Contexts {
		if context.Name == name {
			return context
		}
	}

	return Context{}
}

// GetUser return user by name
func (c *Config) GetUser(name string) User {
	for _, user := range c.Users {
		if user.Name == name {
			return user
		}
	}

	return User{}
}

// GetEndpoint return endpoint by name
func (c *Config) GetEndpoint(name string) Endpoint {
	for _, endpoint := range c.Endpoints {
		if endpoint.Name == name {
			return endpoint
		}
	}

	return Endpoint{}
}
