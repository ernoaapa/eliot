package source

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/layeryd/model"
	"gopkg.in/yaml.v2"
)

// FileSource is source what reads desired state from file
type FileSource struct {
	filePath string
	interval time.Duration
}

// NewFileSource creates new file source what updates the state intervally
func NewFileSource(filePath string, interval time.Duration) *FileSource {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Panicf("Unable to open state, file [%s] does not exist!", filePath)
	}
	return &FileSource{
		filePath,
		interval,
	}
}

// GetUpdates return channel for state changes
func (s *FileSource) GetUpdates(info model.DeviceInfo) chan []model.Pod {
	updates := make(chan []model.Pod)
	go func() {
		for {
			pods, err := s.getPods(info)

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

func (s *FileSource) getPods(info model.DeviceInfo) (pods []model.Pod, err error) {
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
	pods := &[]model.Pod{}
	err := yaml.Unmarshal(data, pods)
	return *pods, err
}
