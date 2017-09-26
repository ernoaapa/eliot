package mapping

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
)

// MapModelByPodNamesToInternalModel builds from container label information pod data structure
func MapModelByPodNamesToInternalModel(containers []containerd.Container) map[string][]model.Container {
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

		result[podName] = append(result[podName], MapContainerToInternalModel(container))
	}
	return result
}

// MapContainersToInternalModel maps containerd models to internal model
func MapContainersToInternalModel(containers []containerd.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, MapContainerToInternalModel(container))
	}
	return result
}

// MapContainerToInternalModel maps containerd model to internal model
func MapContainerToInternalModel(container containerd.Container) model.Container {
	labels := ContainerLabels(container.Info().Labels)
	containerName := labels.getContainerName()
	return model.Container{
		ID:    container.ID(),
		Name:  containerName,
		Image: container.Info().Image,
	}
}
