package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/layeryd/model"
	"gopkg.in/yaml.v2"
)

// FileManifestSource is source what reads manifest from file
type FileManifestSource struct {
	filePath string
	interval time.Duration
}

// NewFileManifestSource creates new file source what updates the state intervally
func NewFileManifestSource(filePath string, interval time.Duration) *FileManifestSource {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Panicf("Unable to open state, file [%s] does not exist!", filePath)
	}
	return &FileManifestSource{
		filePath,
		interval,
	}
}

// GetUpdates return channel for state changes
func (s *FileManifestSource) GetUpdates() chan []model.Pod {
	updates := make(chan []model.Pod)
	go func() {
		for {
			pods, err := s.getPods()

			if err != nil {
				log.Printf("Error reading state: %s", err)
			} else {
				updates <- pods
			}
			time.Sleep(s.interval)
		}
	}()
	return updates
}

func (s *FileManifestSource) getPods() (pods []model.Pod, err error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return pods, fmt.Errorf("Cannot update state, file [%s] does not exist", s.filePath)
		}
		return pods, err
	}

	switch extension := filepath.Ext(s.filePath); extension {
	case ".yaml", ".yml":
		return unmarshalYaml(data)
	default:
		return pods, fmt.Errorf("Invalid source file format: %s", extension)
	}
}

func unmarshalYaml(data []byte) ([]model.Pod, error) {
	target := &[]model.Pod{}

	unmarshalErr := yaml.Unmarshal(data, target)
	if unmarshalErr != nil {
		return []model.Pod{}, errors.Wrapf(unmarshalErr, "Unable to read yaml file")
	}

	pods := model.Defaults(*target)

	validationErr := model.Validate(pods)
	if validationErr != nil {
		return pods, errors.Wrapf(validationErr, "Invalid pod definitions")
	}

	return pods, nil
}
