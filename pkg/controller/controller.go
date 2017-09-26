package controller

import (
	"time"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Controller is responsible for keeping the containers in desired state
type Controller struct {
	client   runtime.Client
	interval time.Duration
	in       <-chan []model.Pod
	out      chan<- []model.Pod
	manifest podsManifest
}

// New creates new container controller
func New(client runtime.Client, interval time.Duration, in <-chan []model.Pod, out chan<- []model.Pod) *Controller {
	return &Controller{
		client:   client,
		interval: interval,
		in:       in,
		out:      out,
	}
}

// Start the controller sync process
// Waits for update from in channel or after given interval runs sync()
func (c *Controller) Start() {
	for {
		select {
		case update := <-c.in:
			c.manifest = update
			err := c.Sync(c.manifest)
			if err != nil {
				log.Warnf("Failed to update container state: %s", err)
			}
		case <-time.After(c.interval):
			err := c.Sync(c.manifest)
			if err != nil {
				log.Warnf("Failed to update container state: %s", err)
			}
		}
	}
}

// Sync start and stop containers to match with target pods
func (c *Controller) Sync(manifest podsManifest) (err error) {
	log.Debugf("Sync containers: %v", manifest)
	namespaces, err := c.client.GetNamespaces()
	if err != nil {
		return errors.Wrapf(err, "Failed to list namespaces when syncing containers")
	}
	namespaces = utils.MergeLists(namespaces, manifest.getNamespaces())

	log.Debugf("Syncing namespaces: %s", namespaces)
	for _, namespace := range namespaces {
		manifest := manifest.filterPodsByNamespace(namespace)
		state, err := c.client.GetAllContainers(namespace)
		if err != nil {
			return err
		}

		if err := c.cleanupRemovedContainers(namespace, manifest, state); err != nil {
			return err
		}

		if err := c.createMissingContainers(namespace, manifest, state); err != nil {
			return err
		}

		if err := c.ensureContainerTasksRunning(manifest, state); err != nil {
			return err
		}
	}

	state, err := getCurrentState(c.client)
	if err != nil {
		return errors.Wrapf(err, "Failed to forward current state, failed to resolve current state!")
	}

	log.Debugf("Controller sync completed!")

	select {
	case c.out <- state:
		return nil
	default:
		log.Warnf("Controller state output is blocking, state reporter not processing messages?")
		return nil
	}
}

func (c *Controller) cleanupRemovedContainers(namespace string, pods podsManifest, state containersState) error {
	remove := getRemovedContainers(pods, state)

	if len(remove) > 0 {
		log.WithFields(log.Fields{
			"namespace": namespace,
			"remove":    len(remove),
		}).Debugf("Remove containers from namespace %s", namespace)

		for _, container := range remove {
			err := c.client.StopContainer(container.ID)
			if err != nil {
				return err
			}
		}
	} else {
		log.Debugf("No containers to remove from namespace %s", namespace)
	}
	return nil
}

func getRemovedContainers(pods podsManifest, state containersState) (remove []model.Container) {
	for podName, containers := range state {
		if !pods.containsPod(podName) {
			log.Debugf("Found active pod [%s], but it does not exist in manifest, will remove it", podName)
			remove = append(remove, containers...)
		} else {
			for _, activeContainer := range containers {
				if !pods.containsContainer(podName, activeContainer) {
					log.Debugf("Found from pod [%s] a container [%s] but it does not exist in manifest, will remove it", podName, activeContainer.Name)
					remove = append(remove, activeContainer)
				}
			}
		}
	}

	return remove
}

func (c *Controller) createMissingContainers(namespace string, pods podsManifest, state containersState) error {
	for _, pod := range pods {
		create := getMissingContainers(pod, state)

		if len(create) > 0 {
			log.WithFields(log.Fields{
				"namespace": namespace,
				"create":    len(create),
			}).Debugf("Missing containers in namespace %s", namespace)

			for _, container := range create {
				createErr := c.client.CreateContainer(pod, container)
				if createErr != nil {
					return errors.Wrapf(createErr, "Failed to create container %s %s", pod.Metadata.Name, container.Name)
				}
			}
		} else {
			log.Debugf("No missing containers in namespace %s for pod %s", namespace, pod.Metadata.Name)
		}
	}
	return nil
}

func getMissingContainers(pod model.Pod, state containersState) (create []model.Container) {
	for _, desiredContainer := range pod.Spec.Containers {
		if !state.containsContainer(pod.Metadata.Name, desiredContainer) {
			create = append(create, desiredContainer)
		}
	}
	return create
}

func (c *Controller) ensureContainerTasksRunning(pods podsManifest, state containersState) error {
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if state.containsContainer(pod.Metadata.Name, container) {
				existingContainerID := state.findContainerID(pod.Metadata.Name, container)
				running, err := c.client.IsContainerRunning(existingContainerID)
				if err != nil {
					return errors.Wrapf(err, "Cannot ensure existing container task running state, get container task returned unexpected error")
				}
				if !running {
					log.Warnf("Detected existing container not running, restarting container [%s]", existingContainerID)
					startErr := c.client.StartContainer(existingContainerID)
					if startErr != nil {
						return startErr
					}
				} else {
					log.Debugf("Container [%s] running and healthy", existingContainerID)
				}
			} else {
				startErr := c.client.StartContainer(container.ID)
				if startErr != nil {
					return errors.Wrapf(startErr, "Error while starting new container %s", container.ID)
				}
			}
		}
	}
	return nil
}
