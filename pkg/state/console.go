package state

import (
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	log "github.com/sirupsen/logrus"
)

// ConsoleStateReporter is Reporter implementation what just prints status to stdout
type ConsoleStateReporter struct {
	resolver *device.Resolver
	in       <-chan []model.Pod
}

// NewConsoleStateReporter creates new ConsoleStateReporter
func NewConsoleStateReporter(resolver *device.Resolver, in <-chan []model.Pod) *ConsoleStateReporter {
	return &ConsoleStateReporter{
		resolver,
		in,
	}
}

// Start printing state to console
func (r *ConsoleStateReporter) Start() {
	for {
		r.report(<-r.in)
	}
}

// Report implements Reporter interface by printing out the state to console
func (r *ConsoleStateReporter) report(state []model.Pod) error {
	for _, pod := range state {
		states := getContainerStateCounts(pod.Status.ContainerStatuses)
		log.WithFields(states).Infof("%s pod containers state", pod.Metadata.Name)
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
