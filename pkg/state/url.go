package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// URLStateReporter is Reporter implementation what PUT data to url
type URLStateReporter struct {
	resolver *device.Resolver
	client   runtime.Client
	interval time.Duration
	url      string
}

// NewURLStateReporter creates new URLStateReporter
func NewURLStateReporter(resolver *device.Resolver, client runtime.Client, interval time.Duration, url string) *URLStateReporter {
	return &URLStateReporter{
		resolver,
		client,
		interval,
		url,
	}
}

// Start starts printing status to console with given interval
func (r *URLStateReporter) Start() {
	for {
		states, err := getCurrentState(r.client)
		if err != nil {
			log.Errorf("Error while reporting current device state: %s", err)
		} else {
			r.report(states)
		}
		time.Sleep(r.interval)
	}
}

// Report implements Reporter interface by printing out the state to console
func (r *URLStateReporter) report(podsWithStates []*model.Pod) error {
	body, err := json.Marshal(podsWithStates)
	if err != nil {
		return errors.Wrap(err, "Error while marshalling device info to JSON")
	}

	resp, err := put(r.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrapf(err, "Cannot PUT current state to [%s]", r.url)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		log.Debugf("Received error response (code %d): %s", resp.StatusCode, string(data[:]))
		return fmt.Errorf("Url replied with status code [%d]", resp.StatusCode)
	}
	return nil
}

func put(url, contentType string, data io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
