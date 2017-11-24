package config

import (
	"io/ioutil"
	"log"

	"github.com/ernoaapa/eliot/pkg/fs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ProjectConfig represents configuration in project directory
type ProjectConfig struct {
	path  string
	Name  string     `yaml:"name"`
	Image string     `yaml:"image"`
	Env   []string   `yaml:"env"`
	Binds []string   `yaml:"binds"`
	Sync  SyncConfig `yaml:"sync"`
}

type SyncConfig struct {
	Image  string `yaml:"image"`
	Target string `yaml:"target"`
}

// NameOrElse return 'name' from the project configuration or defaults to given value
func (p ProjectConfig) NameOrElse(val string) string {
	if p.Name != "" {
		return p.Name
	}
	return val
}

// ImageOrElse return 'name' from the project configuration or defaults to given value
func (p ProjectConfig) ImageOrElse(val string) string {
	if p.Image != "" {
		return p.Image
	}
	return val
}

//
func (p ProjectConfig) EnvWith(values []string) (result []string) {
	result = append(result, p.Env...)
	result = append(result, values...)
	return result
}

// ReadProjectConfig read ProjectConfig from given path
// If file doesn't exist, returns empty config so it's safe for reading even
// The file doesn't exist. In any other failure case, will fatal.
func ReadProjectConfig(path string) (result *ProjectConfig) {
	// Defaults
	config := &ProjectConfig{
		path: path,
		Sync: SyncConfig{
			Image:  "docker.io/stefda/rsync:latest",
			Target: "/volume",
		},
	}

	if !fs.FileExist(path) {
		return config
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(errors.Wrapf(err, "Failed to read project configuration at %s", path))
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln(errors.Wrapf(err, "Error while reading YAML configuration at %s", path))
	}

	return config
}
