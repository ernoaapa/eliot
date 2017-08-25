package state

import (
	"fmt"
	"time"

	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/runtime"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
)

// ConsoleStateReporter is Reporter implementation what just prints status to stdout
type ConsoleStateReporter struct {
	info     *model.DeviceInfo
	client   *runtime.ContainerdClient
	interval time.Duration
}

// NewConsoleStateReporter creates new ConsoleStateReporter
func NewConsoleStateReporter(info *model.DeviceInfo, client *runtime.ContainerdClient, interval time.Duration) *ConsoleStateReporter {
	return &ConsoleStateReporter{
		info,
		client,
		interval,
	}
}

// Start starts printing status to console with given interval
func (r *ConsoleStateReporter) Start() {
	r.registerDevice()
	r.runReportLoop()
}

func (r *ConsoleStateReporter) registerDevice() {
	log.WithFields(structs.Map(r.info)).Infoln("Device registered")
}

func (r *ConsoleStateReporter) runReportLoop() {
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
func (r *ConsoleStateReporter) report(states map[string]*model.DeviceState) error {

	for namespace, state := range states {
		log.WithFields(log.Fields{
			"nr of pods": fmt.Sprintf("%d containers", len(state.Pods)),
		}).Infof("%s state update", namespace)
	}
	return nil
}
