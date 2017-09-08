package config

import (
	"io/ioutil"

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
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

// Context represents single user in single endpoint
type Context struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Endpoint string `yaml:"endpoint"`
}

// GetConfig reads current config from user home directory
func GetConfig(path string) (*Config, error) {
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

// GetCurrentUser return current user
func (c *Config) GetCurrentUser() User {
	return c.GetUser(c.GetCurrentContext().User)
}

// GetCurrentContext return current context
func (c *Config) GetCurrentContext() Context {
	return c.GetContext(c.CurrentContext)
}

// GetCurrentEndpoint return current context
func (c *Config) GetCurrentEndpoint() Endpoint {
	return c.GetEndpoint(c.GetCurrentContext().Endpoint)
}

// GetContext return context by name
func (c *Config) GetContext(name string) Context {
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
