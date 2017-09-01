package manifest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
)

// URLManifestSource is source what reads manifest from file
type URLManifestSource struct {
	manifestURL string
	interval    time.Duration
	resolver    *device.Resolver
}

// NewURLManifestSource creates new url source what updates the state intervally
func NewURLManifestSource(manifestURL string, interval time.Duration, resolver *device.Resolver) *URLManifestSource {
	return &URLManifestSource{
		manifestURL: manifestURL,
		interval:    interval,
		resolver:    resolver,
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
				log.Warnf("Error while fetching manifest: %s", err)
			} else {
				updates <- pods
			}
			time.Sleep(s.interval)
		}
	}()
	return updates
}

func (s *URLManifestSource) getPods() (pods []model.Pod, err error) {
	body, err := json.Marshal(s.resolver.GetInfo())
	if err != nil {
		return pods, errors.Wrap(err, "Error wile marshalling device info to JSON")
	}
	resp, err := http.Post(s.manifestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return pods, errors.Wrapf(err, "Cannot download manifest file [%s]", s.manifestURL)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pods, errors.Wrapf(err, "Failed to read response")
	}

	if strings.HasSuffix(strings.ToLower(s.manifestURL), "json") {
		return unmarshalJSON(data)
	}
	return unmarshalYaml(data)
}
