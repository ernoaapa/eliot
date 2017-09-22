package controller

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Sync start and stop containers to match with target pods
func Sync(client runtime.Client, allPods []model.Pod) (err error) {
	log.Debugf("Received update, start updating containerd: %v", allPods)
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return errors.Wrapf(err, "Failed to list namespaces when syncing containers")
	}
	namespaces = utils.MergeLists(namespaces, getNamespaces(allPods))

	log.Debugf("Found namespaces: %s", namespaces)
	for _, namespace := range namespaces {
		pods := filterByNamespace(namespace, allPods)
		containers, err := client.GetContainers(namespace)
		if err != nil {
			return err
		}

		if len(containers) > 0 {
			log.Debugf("Found %d containers in namespace %s", len(containers), namespace)
		} else {
			log.Debugf("No running containers in namespace %s", namespace)
		}

		if err := cleanupRemovedContainers(client, namespace, pods, containers); err != nil {
			return err
		}

		active, err := createMissingContainers(client, namespace, pods, containers)
		if err != nil {
			return err
		}

		if err := ensureContainerTasksRunning(client, namespace, pods, active); err != nil {
			return err
		}
	}
	return nil
}

func cleanupRemovedContainers(client runtime.Client, namespace string, pods []model.Pod, containers []containerd.Container) error {

	remove := getRemovedContainers(containers, pods)

	if len(remove) > 0 {
		log.WithFields(log.Fields{
			"namespace": namespace,
			"remove":    len(remove),
		}).Debugf("Remove containers from namespace %s", namespace)

		for _, container := range remove {
			err := client.StopContainer(container)
			if err != nil {
				return err
			}
		}
	} else {
		log.Debugf("No containers to remove from namespace %s", namespace)
	}
	return nil
}

func createMissingContainers(client runtime.Client, namespace string, pods []model.Pod, containers []containerd.Container) ([]containerd.Container, error) {
	createdContainers := []containerd.Container{}
	for _, pod := range pods {
		create := getMissingContainers(pod, containers)

		if len(create) > 0 {
			log.WithFields(log.Fields{
				"namespace": namespace,
				"create":    len(create),
			}).Debugf("Missing containers in namespace %s", namespace)

			for _, container := range create {
				created, createErr := client.CreateContainer(pod, container)
				if createErr != nil {
					return nil, errors.Wrapf(createErr, "Failed to create container %s %s", pod.GetName(), container.Name)
				}
				createdContainers = append(createdContainers, created)
			}
		} else {
			log.Debugf("No missing containers in namespace %s for pod %s", namespace, pod.GetName())
		}
	}
	return append(containers, createdContainers...), nil
}

func ensureContainerTasksRunning(client runtime.Client, namespace string, pods []model.Pod, containers []containerd.Container) error {
	for _, container := range containers {
		running, err := client.IsContainerRunning(container)
		if err != nil {
			return errors.Wrapf(err, "Cannot ensure container task running, get container task returned unexpected error")
		}
		if !running {
			startErr := client.StartContainer(container)
			if startErr != nil {
				return startErr
			}
		}
	}
	return nil
}

func getNamespaces(pods []model.Pod) []string {
	result := []string{}
	for _, pod := range pods {
		namespace := pod.GetNamespace()
		if namespace != "" {
			result = append(result, pod.GetNamespace())
		}
	}
	return result
}

func getRemovedContainers(active []containerd.Container, pods []model.Pod) (remove []containerd.Container) {
	for _, activeContainer := range active {
		if !containsTargetContainer(activeContainer, pods) {
			remove = append(remove, activeContainer)
		}
	}

	return
}

func containsTargetContainer(target containerd.Container, pods []model.Pod) bool {
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if containerMatch(pod, container, target) {
				return true
			}
		}
	}
	return false
}

func containerMatch(pod model.Pod, container model.Container, target containerd.Container) bool {
	labels := target.Info().Labels
	podName := runtime.GetLabelFor(labels, runtime.PodNameSuffix)
	podNamespace := runtime.GetLabelFor(labels, runtime.PodNamespaceSuffix)
	containerName := runtime.GetLabelFor(labels, runtime.ContainerNameSuffix)

	return pod.GetNamespace() == podNamespace &&
		pod.GetName() == podName &&
		container.Name == containerName
}

func getMissingContainers(pod model.Pod, active []containerd.Container) (create []model.Container) {
	for _, targetContainer := range pod.Spec.Containers {
		if !containsActiveContainer(pod, targetContainer, active) {
			create = append(create, targetContainer)
		}
	}
	return
}

func containsActiveContainer(pod model.Pod, target model.Container, list []containerd.Container) bool {
	for _, item := range list {
		if containerMatch(pod, target, item) {
			return true
		}
	}
	return false
}

func isUpToDate(ctx context.Context, target model.Container, active containerd.Container) bool {
	if target.Image != active.Info().Image {
		return false
	}

	if !isContainerRunning(ctx, active) {
		return false
	}

	return true
}

func isContainerRunning(ctx context.Context, container containerd.Container) bool {
	task, err := container.Task(ctx, nil)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false
		}
		log.Warnf("Unable to get container task, assuming it's running. Error was: %s", err)
		return true
	}

	status, err := task.Status(ctx)
	if err != nil {
		log.Warnf("Unable to get container task state, assuming it's running. Error was: %s", err)
		return true
	}
	return status.Status == containerd.Running
}

func filterByNamespace(namespace string, pods []model.Pod) []model.Pod {
	result := []model.Pod{}
	for _, pod := range pods {
		if namespace == pod.GetNamespace() {
			result = append(result, pod)
		}
	}
	return result
}
