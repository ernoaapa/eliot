package containerd

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
)

// MapToModelByPodNames builds from container label information pod data structure
func MapToModelByPodNames(containers []containerd.Container) map[string][]model.Container {
	result := make(map[string][]model.Container)
	for _, container := range containers {
		labels := ContainerLabels(container.Info().Labels)
		podName := labels.getPodName()
		if podName == "" {
			// container is not cand managed container so add it under 'system' pod in namespace 'default'
			podName = "system"
		}

		if _, ok := result[podName]; !ok {
			result[podName] = []model.Container{}
		}

		result[podName] = append(result[podName], MapContainerToModel(container))
	}
	return result
}

// MapContainersToModel maps containerd models to internal model
func MapContainersToModel(containers []containerd.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, MapContainerToModel(container))
	}
	return result
}

// MapContainerToModel maps containerd model to internal model
func MapContainerToModel(container containerd.Container) model.Container {
	labels := ContainerLabels(container.Info().Labels)
	containerName := labels.getContainerName()
	return model.Container{
		ID:    container.ID(),
		Name:  containerName,
		Image: container.Info().Image,
	}
}
