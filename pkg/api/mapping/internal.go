package mapping

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

// MapPodsToAPIModel maps list of internal Pod models to API model
func MapPodsToAPIModel(namespace string, containersByPods map[string][]model.Container) (result []*pb.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, MapPodToAPIModel(namespace, podName, containers))
	}
	return result
}

// MapPodToAPIModel maps internal Pod model to API model
func MapPodToAPIModel(namespace, podName string, containers []model.Container) *pb.Pod {
	return &pb.Pod{
		Metadata: &pb.ResourceMetadata{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: &pb.PodSpec{
			Containers: MapContainersToAPIModel(containers),
		},
		// Status: &pb.PodStatus{
		// 	ContainerStatuses: mapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
		// },
	}
}

// MapContainersToAPIModel maps list of internal Container models to API model
func MapContainersToAPIModel(containers []model.Container) (result []*pb.Container) {
	for _, container := range containers {
		result = append(result, &pb.Container{
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}

// MapContainerStatusesToAPIModel maps list of internal ContainerStatus models to API model
func MapContainerStatusesToAPIModel(statuses []model.ContainerStatus) (result []*pb.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &pb.ContainerStatus{
			ContainerID: status.ContainerID,
			Image:       status.Image,
			State:       status.State,
		})
	}
	return result
}
