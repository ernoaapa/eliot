package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// URLStateReporter is Reporter implementation what PUT data to url
type URLStateReporter struct {
	resolver *device.Resolver
	in       <-chan []model.Pod
	url      string
}

// NewURLStateReporter creates new URLStateReporter
func NewURLStateReporter(resolver *device.Resolver, in <-chan []model.Pod, url string) *URLStateReporter {
	return &URLStateReporter{
		resolver,
		in,
		url,
	}
}

// Start reporting state to the given url
func (r *URLStateReporter) Start() {
	for {
		r.report(<-r.in)
	}
}

// Report implements Reporter interface by printing out the state to console
func (r *URLStateReporter) report(state []model.Pod) error {
	log.Debugf("Received updated state, send it to url [%s]", r.url)

	body, err := json.Marshal(state)
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

	log.Debugf("Sent state report succesfully!")
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
