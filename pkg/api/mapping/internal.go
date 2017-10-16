package mapping

import (
	core "github.com/ernoaapa/can/pkg/api/core"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

// MapPodsToAPIModel maps list of internal pod models to API model
func MapPodsToAPIModel(pods []model.Pod) (result []*pods.Pod) {
	for _, pod := range pods {
		result = append(result, MapPodToAPIModel(pod))
	}
	return result
}

// MapContainersByPodsToAPIModel maps list of internal Pod models to API model
func MapContainersByPodsToAPIModel(namespace string, containersByPods map[string][]model.Container) (result []*pods.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, CreatePodAPIModel(podName, namespace, containers))
	}
	return result
}

// CreatePodAPIModel maps internal Pod model to API model
func CreatePodAPIModel(namespace, podName string, containers []model.Container) *pods.Pod {
	pod := model.NewPod(podName, namespace)
	pod.Spec.Containers = containers
	return MapPodToAPIModel(pod)
}

// MapPodToAPIModel maps internal Pod model to API model
func MapPodToAPIModel(pod model.Pod) *pods.Pod {
	return &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
		},
		Spec: &pods.PodSpec{
			Containers: MapContainersToAPIModel(pod.Spec.Containers),
		},
		Status: &pods.PodStatus{
			ContainerStatuses: MapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
		},
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
