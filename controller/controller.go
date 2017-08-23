package controller

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/ernoaapa/layeryd/model"
	log "github.com/sirupsen/logrus"
)

// Sync start and stop containers to match with target pods
func Sync(ctx context.Context, client *containerd.Client, pods []model.Pod) error {
	log.Debugln("Received update, start updating containerd")

	containers, err := client.Containers(ctx)
	if err != nil {
		log.Warnf("Error getting list of containers: %v", err)
		return err
	}
	for _, pod := range pods {
		create, remove := groupContainers(ctx, pod, containers)

		log.WithFields(log.Fields{
			"running": len(containers),
			"create":  len(create),
			"remove":  len(remove),
		}).Debugln("Resolved current container status")

		if err := createContainers(ctx, client, pod, create); err != nil {
			return err
		}
		if err := stopContainers(ctx, remove); err != nil {
			return err
		}
	}
	return nil
}

func groupContainers(ctx context.Context, pod model.Pod, active []containerd.Container) (create []model.Container, remove []containerd.Container) {

	for _, targetContainer := range pod.Spec.Containers {
		if !containsActiveContainer(pod.GetName(), targetContainer, active) {
			create = append(create, targetContainer)
		}
	}

	for _, activeContainer := range active {
		if !containsTargetContainer(pod.GetName(), activeContainer, pod.Spec.Containers) {
			remove = append(remove, activeContainer)
		}
	}

	return
}

func containsTargetContainer(podName string, target containerd.Container, list []model.Container) bool {
	for _, item := range list {
		if target.ID() == item.BuildID(podName) {
			return true
		}
	}
	return false
}

func containsActiveContainer(podName string, target model.Container, list []containerd.Container) bool {
	for _, item := range list {
		if target.BuildID(podName) == item.ID() {
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
