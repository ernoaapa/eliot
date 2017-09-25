package api

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

func mapPodsToAPIModel(namespace string, containersByPods map[string][]model.Container) (result []*pb.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, &pb.Pod{
			Metadata: model.NewMetadata(podName, namespace),
			Spec: &pb.PodSpec{
				Containers: mapContainersToAPIModel(containers),
			},
			// Status: &pb.PodStatus{
			// 	ContainerStatuses: mapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
			// },
		})
	}
	return result
}

func mapContainersToAPIModel(containers []model.Container) (result []*pb.Container) {
	for _, container := range containers {
		result = append(result, &pb.Container{
			ID:    container.ID,
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}

func mapContainerStatusesToAPIModel(statuses []model.ContainerStatus) (result []*pb.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &pb.ContainerStatus{
			ContainerID: status.ContainerID,
			Image:       status.Image,
			State:       status.State,
		})
	}
	return result
}
