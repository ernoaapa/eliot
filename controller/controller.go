package controller

import (
	"context"

	"github.com/containerd/containerd"
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

	remove, create := groupContainers(target.Containers, getIds(containers))
	log.Printf("Remove %d containers", len(remove))
	log.Printf("Create %d containers", len(create))
	return nil
}

func groupContainers(target, active []string) (remove []string, create []string) {
	for _, targetContainer := range target {
		if !contains(targetContainer, active) {
			create = append(create, targetContainer)
		}
	}

	for _, activeContainer := range active {
		if !contains(activeContainer, target) {
			remove = append(remove, activeContainer)
		}
	}

	return
}

func getIds(containers []containerd.Container) []string {
	result := []string{}
	for _, container := range containers {
		result = append(result, container.ID())
	}
	return result
}

func contains(target string, list []string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}
