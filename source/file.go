package source

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ernoaapa/layeryd/model"
	"gopkg.in/yaml.v2"
)

type FileSource struct {
	filePath string
	interval time.Duration
}

func NewFileSource(filePath string, interval time.Duration) *FileSource {
	return &FileSource{
		filePath,
		interval,
	}
}

func (s *FileSource) GetUpdates(info model.NodeInfo) chan model.Pod {
	updates := make(chan model.Pod)
	go func() {
		for {
			state, err := s.GetState(info)

			if err != nil {
				log.Printf("Error reading state: %s", err)
			} else {
				updates <- state
			}
			time.Sleep(s.interval)
		}
	}()
	return updates
}

func (s *FileSource) GetState(info model.NodeInfo) (pod model.Pod, err error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return pod, fmt.Errorf("Unable to open state, file [%s] does not exist!", s.filePath)
		}
		return pod, err
	}

	switch extension := filepath.Ext(s.filePath); extension {
	case ".yaml", ".yml":
		return unmarshalYaml(data)
	default:
		return pod, fmt.Errorf("Invalid source file format: %s", extension)
	}
}

func unmarshalYaml(data []byte) (model.Pod, error) {
	pod := &model.Pod{}
	err := yaml.Unmarshal(data, pod)
	return *pod, err
}
