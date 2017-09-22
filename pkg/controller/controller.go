package controller

import (
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Sync start and stop containers to match with target pods
func Sync(client runtime.Client, manifest podsManifest) (err error) {
	log.Debugf("Received update, start updating containerd: %v", manifest)
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return errors.Wrapf(err, "Failed to list namespaces when syncing containers")
	}
	namespaces = utils.MergeLists(namespaces, manifest.getNamespaces())

	log.Debugf("Syncing namespaces: %s", namespaces)
	for _, namespace := range namespaces {
		manifest := manifest.filterPodsByNamespace(namespace)
		state, err := client.GetContainersByPods(namespace)
		if err != nil {
			return err
		}

		if err := cleanupRemovedContainers(client, namespace, manifest, state); err != nil {
			return err
		}

		if err := createMissingContainers(client, namespace, manifest, state); err != nil {
			return err
		}

		if err := ensureContainerTasksRunning(client, manifest, state); err != nil {
			return err
		}
	}
	return nil
}

func cleanupRemovedContainers(client runtime.Client, namespace string, pods podsManifest, state containersState) error {
	remove := getRemovedContainers(pods, state)

	if len(remove) > 0 {
		log.WithFields(log.Fields{
			"namespace": namespace,
			"remove":    len(remove),
		}).Debugf("Remove containers from namespace %s", namespace)

		for _, container := range remove {
			err := client.StopContainer(container.ID)
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

func createMissingContainers(client runtime.Client, namespace string, pods podsManifest, state containersState) error {
	for _, pod := range pods {
		create := getMissingContainers(pod, state)

		if len(create) > 0 {
			log.WithFields(log.Fields{
				"namespace": namespace,
				"create":    len(create),
			}).Debugf("Missing containers in namespace %s", namespace)

			for _, container := range create {
				createErr := client.CreateContainer(pod, container)
				if createErr != nil {
					return errors.Wrapf(createErr, "Failed to create container %s %s", pod.GetName(), container.Name)
				}
			}
		} else {
			log.Debugf("No missing containers in namespace %s for pod %s", namespace, pod.GetName())
		}
	}
	return nil
}

func getMissingContainers(pod model.Pod, state containersState) (create []model.Container) {
	for _, desiredContainer := range pod.Spec.Containers {
		if !state.containsContainer(pod.GetName(), desiredContainer) {
			create = append(create, desiredContainer)
		}
	}
	return create
}

func ensureContainerTasksRunning(client runtime.Client, pods podsManifest, state containersState) error {
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if state.containsContainer(pod.GetName(), container) {
				existingContainerID := state.findContainerID(pod.GetName(), container)
				running, err := client.IsContainerRunning(existingContainerID)
				if err != nil {
					return errors.Wrapf(err, "Cannot ensure existing container task running state, get container task returned unexpected error")
				}
				if !running {
					log.Warnf("Detected existing container not running, restarting container [%s]", existingContainerID)
					startErr := client.StartContainer(existingContainerID)
					if startErr != nil {
						return startErr
					}
				} else {
					log.Debugf("Container [%s] running and healthy", existingContainerID)
				}
			} else {
				startErr := client.StartContainer(container.ID)
				if startErr != nil {
					return errors.Wrapf(startErr, "Error while starting new container %s", container.ID)
				}
			}
		}
	}
	return nil
}
