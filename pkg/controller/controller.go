package controller

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Sync start and stop containers to match with target pods
func Sync(client *runtime.ContainerdClient, allPods []model.Pod) (err error) {
	log.Debugln("Received update, start updating containerd")
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return errors.Wrapf(err, "Failed to list namespaces when syncing containers")
	}

	log.Debugf("Found namespaces: %s", namespaces)
	for _, namespace := range namespaces {
		pods := filterByNamespace(namespace, allPods)
		containers, err := client.GetContainers(namespace)
		log.Debugf("Found running containers in namespace %s: %d", namespace, len(containers))
		if err != nil {
			return err
		}

		for _, pod := range pods {
			create := getMissingContainers(pod, containers)

			log.WithFields(log.Fields{
				"namespace": namespace,
				"create":    len(create),
				"running":   len(containers),
			}).Debugf("Missing containers in namespace %s", namespace)

			if err := client.CreateContainers(pod, create); err != nil {
				return err
			}
		}

		remove := getRemovedContainers(containers, pods)

		log.WithFields(log.Fields{
			"namespace": namespace,
			"remove":    len(remove),
		}).Debugf("Remove containers from namespace %s", namespace)

		if err := client.StopContainers(remove); err != nil {
			return err
		}
	}
	return nil
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
			if target.ID() == container.ID {
				return true
			}
		}
	}
	return false
}

func getMissingContainers(pod model.Pod, active []containerd.Container) (create []model.Container) {
	for _, targetContainer := range pod.Spec.Containers {
		if !containsActiveContainer(targetContainer, active) {
			create = append(create, targetContainer)
		}
	}
	return
}

func containsActiveContainer(target model.Container, list []containerd.Container) bool {
	for _, item := range list {
		if target.ID == item.ID() {
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
