package controller

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/ernoaapa/eliot/pkg/runtime"
	log "github.com/sirupsen/logrus"
)

// Lifecycle is controller which monitors containers and if container stops,
// restart it based on restart policy
type Lifecycle struct {
	client   runtime.Client
	interval time.Duration
}

// NewLifecycle creates new Lifecycle controller instance
func NewLifecycle(client runtime.Client) *Lifecycle {
	return &Lifecycle{
		client:   client,
		interval: 5 * time.Second,
	}
}

// Run starts the controller to monitor containers
// Run until faces such fatal error that cannot recover and will return the error
func (l *Lifecycle) Run() error {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for range ticker.C {
		err := l.checkAll()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Lifecycle) checkAll() error {
	namespaces, err := l.client.GetNamespaces()
	if err != nil {
		log.Warnf("Lifecycle controller cannot validate container statuses, error while fetching namespaces: %s", err)
		return nil
	}

	for _, namespace := range namespaces {
		pods, err := l.client.GetPods(namespace)
		if err != nil {
			log.Warnf("Lifecycle controller cannot validate container statuses, error while fetching pods: %s", err)
			continue
		}

		for _, pod := range pods {
			for _, status := range pod.Status.ContainerStatuses {
				if status.State == "stopped" || status.State == "unknown" && pod.Spec.RestartPolicy == "always" {
					log.Debugf("Detected [%s] container [%s] in namespace [%s] with 'always' restart policy", status.State, status.ContainerID, pod.Metadata.Name)
					ioset, err := runtime.NewIOSet(fmt.Sprintf("%s.%s", pod.Metadata.Name, status.Name))
					if err != nil {
						return errors.Wrapf(err, "Error while creating container ioset, cannot run lifecycle controller")
					}
					status, err := l.client.StartContainer(namespace, status.ContainerID, *ioset)
					if err != nil {
						log.Warnf("Lifecycle controller failed to start container: %s", err)
						continue
					}
					log.Debugf("Restarted container [%s] in namespace [%s]", status.ContainerID, pod.Metadata.Name)
				}
			}
		}
	}
	return nil
}
