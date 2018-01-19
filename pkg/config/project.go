package config

import (
	"io/ioutil"

	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/eliot/pkg/fs"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ProjectConfig represents configuration in project directory
type ProjectConfig struct {
	path    string
	Name    string   `yaml:"name,omitempty"`
	Image   string   `yaml:"image,omitempty"`
	Command []string `yaml:"command,omitempty"`
	Env     []string `yaml:"env,omitempty"`
	Binds   []string `yaml:"binds,omitempty"`
	Mounts  []string `yaml:"mounts,omitempty"`
	WorkDir string   `yaml:"workdir,omitempty"`

	SyncContainer *containers.Container `yaml:"syncContainer,omitempty"`
	Syncs         []string              `yaml:"syncs,omitempty"`
}

// EnvWith return list of environment variable definitions from project configs
// with values appended to end of the list.
func (p ProjectConfig) EnvWith(values []string) (result []string) {
	result = append(result, p.Env...)

	for _, value := range values {
		if !model.IsValidEnvKeyValuePair(value) {
			log.Fatalf("Invalid environment variable [%s], must be in format KEY=value. E.g. --env FOO=bar", value)
		}
		result = append(result, value)
	}
	return result
}

// ReadProjectConfig read ProjectConfig from given path
// If file doesn't exist, returns empty config so it's safe for reading even
// The file doesn't exist. In any other failure case, will fatal.
func ReadProjectConfig(path string) *ProjectConfig {
	// Defaults
	config := &ProjectConfig{
		path: path,
		SyncContainer: &containers.Container{
			Name:  "sync",
			Image: "docker.io/ernoaapa/rsync:ff1a0fa9bc12d051bae4300508d361cc65df338b",
			Env: []string{
				"USER=root",
				"GROUP=root",
			},
		},
		Syncs: []string{
			".:/volume",
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

// WriteProjectConfig writes project config in yaml format into given path
func WriteProjectConfig(path string, config *ProjectConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "Unable to marshal config to yaml format")
	}
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return errors.Wrapf(err, "Failed to write config to path [%s]", path)
	}

	return nil
}
