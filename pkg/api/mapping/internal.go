package mapping

import (
	core "github.com/ernoaapa/can/pkg/api/core"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

// MapPodsToAPIModel maps list of internal Pod models to API model
func MapPodsToAPIModel(namespace string, containersByPods map[string][]model.Container) (result []*pods.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, MapPodToAPIModel(namespace, podName, containers))
	}
	return result
}

// MapPodToAPIModel maps internal Pod model to API model
func MapPodToAPIModel(namespace, podName string, containers []model.Container) *pods.Pod {
	return &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: &pods.PodSpec{
			Containers: MapContainersToAPIModel(containers),
		},
		// Status: &pods.PodStatus{
		// 	ContainerStatuses: mapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
		// },
	}
}

// MapContainersToAPIModel maps list of internal Container models to API model
func MapContainersToAPIModel(source []model.Container) (result []*containers.Container) {
	for _, container := range source {
		result = append(result, &containers.Container{
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}

// MapContainerStatusesToAPIModel maps list of internal ContainerStatus models to API model
func MapContainerStatusesToAPIModel(statuses []model.ContainerStatus) (result []*containers.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &containers.ContainerStatus{
			ContainerID: status.ContainerID,
			Image:       status.Image,
			State:       status.State,
		})
	}
	return result
}
