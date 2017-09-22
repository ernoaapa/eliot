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

		if err := ensureContainerTasksRunning(client, manifest); err != nil {
			return err
		}
	}
	return nil
}

func cleanupRemovedContainers(client runtime.Client, namespace string, pods podsManifest, containers containersState) error {
	remove := getRemovedContainers(pods, containers)

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

func getRemovedContainers(pods podsManifest, containers containersState) (remove []model.Container) {
	for podName, containers := range containers {
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

func createMissingContainers(client runtime.Client, namespace string, pods podsManifest, containers containersState) error {
	for _, pod := range pods {
		create := getMissingContainers(pod, containers)

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

func getMissingContainers(pod model.Pod, containers containersState) (create []model.Container) {
	for _, desiredContainer := range pod.Spec.Containers {
		if !containers.containsContainer(pod.GetName(), desiredContainer) {
			create = append(create, desiredContainer)
		}
	}
	return create
}

func ensureContainerTasksRunning(client runtime.Client, pods podsManifest) error {
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			running, err := client.IsContainerRunning(container.ID)
			if err != nil {
				return errors.Wrapf(err, "Cannot ensure container task running, get container task returned unexpected error")
			}
			if !running {
				startErr := client.StartContainer(container.ID)
				if startErr != nil {
					return startErr
				}
			}
		}
	}
	return nil
}
