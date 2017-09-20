package state

import (
	"time"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	log "github.com/sirupsen/logrus"
)

// ConsoleStateReporter is Reporter implementation what just prints status to stdout
type ConsoleStateReporter struct {
	resolver *device.Resolver
	client   runtime.Client
	interval time.Duration
}

// NewConsoleStateReporter creates new ConsoleStateReporter
func NewConsoleStateReporter(resolver *device.Resolver, client runtime.Client, interval time.Duration) *ConsoleStateReporter {
	return &ConsoleStateReporter{
		resolver,
		client,
		interval,
	}
}

// Start starts printing status to console with given interval
func (r *ConsoleStateReporter) Start() {
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
func (r *ConsoleStateReporter) report(podsWithStates []*model.Pod) error {
	for _, pod := range podsWithStates {
		states := getContainerStateCounts(pod.Status.ContainerStatuses)
		log.WithFields(states).Infof("%s pod containers state", pod.GetName())
	}
	return nil
}

func getContainerStateCounts(statuses []model.ContainerStatus) log.Fields {
	result := map[string]int{}
	for _, status := range statuses {
		if _, ok := result[status.State]; !ok {
			result[status.State] = 0
		}
		result[status.State] = result[status.State] + 1
	}

	fields := log.Fields{}
	for key, value := range result {
		fields[key] = value
	}
	return fields
}
