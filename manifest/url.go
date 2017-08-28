package manifest

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/layeryd/model"
)

// URLManifestSource is source what reads manifest from file
type URLManifestSource struct {
	manifestURL string
	interval    time.Duration
}

// NewURLManifestSource creates new url source what updates the state intervally
func NewURLManifestSource(manifestURL string, interval time.Duration) *URLManifestSource {
	return &URLManifestSource{
		manifestURL,
		interval,
	}
}

// GetUpdates return channel for manifest changes
func (s *URLManifestSource) GetUpdates() chan []model.Pod {
	updates := make(chan []model.Pod)
	go func() {
		for {
			log.Debugf("Load manifest from %s", s.manifestURL)
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

func (s *URLManifestSource) getPods() (pods []model.Pod, err error) {
	resp, err := http.Get(s.manifestURL)
	if err != nil {
		return pods, errors.Wrapf(err, "Cannot download manifest file [%s]", s.manifestURL)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pods, errors.Wrapf(err, "Failed to read response")
	}

	return unmarshalYaml(data)
}
