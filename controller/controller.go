package controller

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/source"
	log "github.com/sirupsen/logrus"
)

func Sync(ctx context.Context, client *containerd.Client, source source.Source) error {

	target, err := source.GetState(model.NodeInfo{})
	if err != nil {
		return err
	}

	containers, err := client.Containers(ctx)
	if err != nil {
		log.Warnf("Error getting list of containers: %v", err)
		return err
	}
	log.Printf("Found %d containers running", len(containers))

	create, valid, update, remove := groupContainers(ctx, target.Spec.Containers, containers)

	if err := createContainers(ctx, client, create); err != nil {
		return err
	}
	if err := stopContainers(ctx, remove); err != nil {
		return err
	}
	log.Printf("Valid containers: %d", len(valid))
	log.Printf("Update containers: %d", len(update))

	return nil
}

func groupContainers(ctx context.Context, target []model.Container, active []containerd.Container) (create []model.Container, valid []containerd.Container, update []containerd.Container, remove []containerd.Container) {
	existing := make(map[model.Container]containerd.Container)
	for _, targetContainer := range target {
		runningContainer := findActiveContainer(targetContainer, active)
		if runningContainer != nil {
			existing[targetContainer] = runningContainer
		} else {
			create = append(create, targetContainer)
		}
	}

	for _, activeContainer := range active {
		if !containsTargetContainer(activeContainer, target) {
			remove = append(remove, activeContainer)
		}
	}

	for targetContainer, existingContainer := range existing {
		if isUpToDate(ctx, targetContainer, existingContainer) {
			valid = append(valid, existingContainer)
		} else {
			update = append(update, existingContainer)
		}
	}

	return
}

func containsTargetContainer(target containerd.Container, list []model.Container) bool {
	for _, item := range list {
		if target.ID() == item.Name {
			return true
		}
	}
	return false
}

func findActiveContainer(target model.Container, list []containerd.Container) containerd.Container {
	for _, item := range list {
		if target.Name == item.ID() {
			return item
		}
	}
	return nil
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
