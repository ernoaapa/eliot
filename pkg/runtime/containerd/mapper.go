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
		containerName := labels.getContainerName()

		if _, ok := result[podName]; !ok {
			result[podName] = []model.Container{}
		}

		result[podName] = append(result[podName], model.Container{
			ID:    container.ID(),
			Name:  containerName,
			Image: container.Info().Image,
		})
	}
	return result
}
