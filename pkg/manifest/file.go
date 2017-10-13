package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/can/pkg/api/mapping"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
)

// FileManifestSource is source what reads manifest from file
type FileManifestSource struct {
	filePath string
	interval time.Duration
	resolver *device.Resolver
	out      chan<- []model.Pod
	running  bool
}

// NewFileManifestSource creates new file source what updates intervally
func NewFileManifestSource(filePath string, interval time.Duration, resolver *device.Resolver, out chan<- []model.Pod) *FileManifestSource {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Panicf("Unable to open state, file [%s] does not exist!", filePath)
	}
	return &FileManifestSource{
		filePath: filePath,
		interval: interval,
		resolver: resolver,
		out:      out,
	}
}

// Start the manifest file update process
func (s *FileManifestSource) Start() {
	s.running = true
	for {
		pods, err := s.getPods()
		if err != nil {
			log.Printf("Error while fetching manifest: %s", err)
		} else {
			s.out <- pods
		}

		time.Sleep(s.interval)
		if !s.running {
			return
		}
	}
}

// Stop the file update polling
func (s *FileManifestSource) Stop() {
	s.running = false
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
		content, err := pb.UnmarshalListYaml(data)
		if err != nil {
			return nil, err
		}
		manifest := mapping.MapPodsToInternalModel(content)
		return manifest, model.Validate(manifest)
	default:
		return pods, fmt.Errorf("Invalid source file format: %s", extension)
	}
}
