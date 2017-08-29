package state

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/layery/pkg/model"
	"github.com/ernoaapa/layery/pkg/runtime"
)

func getCurrentState(client *runtime.ContainerdClient) (result map[string]*model.DeviceState, err error) {
	result = map[string]*model.DeviceState{}

	namespaces, err := client.GetNamespaces()
	if err != nil {
		return result, err
	}

	for _, namespace := range namespaces {
		containers, err := client.GetContainers(namespace)
		if err != nil {
			return result, err
		}

		result[namespace] = &model.DeviceState{
			Pods: mapToPodStates(containers),
		}
	}
	return result, nil
}

func mapToPodStates(containers []containerd.Container) (states []model.PodState) {
	for _, container := range containers {
		states = append(states, model.PodState{
			ID: container.ID(),
		})
	}
	return states
}
