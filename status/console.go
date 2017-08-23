package status

import (
	"fmt"

	"github.com/ernoaapa/layeryd/model"
	log "github.com/sirupsen/logrus"
)

// ConsoleReporter is Reporter implementation what just prints status to stdout
type ConsoleReporter struct {
}

// NewConsoleReporter creates new ConsoleReporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// Report implements Reporter interface by printing out the state to console
func (r *ConsoleReporter) Report(info model.DeviceInfo, state model.DeviceState) error {
	log.WithFields(log.Fields{
		"nr of pods": fmt.Sprintf("%d containers", len(state.Pods)),
	}).Info("State update")
	return nil
}
