package api

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

func mapPodsToApiModel(containersByPods map[string][]model.Container) (result []*pb.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, &pb.Pod{
			Metadata: map[string]string{
				"name": podName,
			},
			Spec: &pb.PodSpec{
				Containers: mapContainersToApiModel(containers),
			},
			// Status: &pb.PodStatus{
			// 	ContainerStatuses: mapContainerStatusesToApiModel(pod.Status.ContainerStatuses),
			// },
		})
	}
	return result
}

func mapContainersToApiModel(containers []model.Container) (result []*pb.Container) {
	for _, container := range containers {
		result = append(result, &pb.Container{
			ID:    container.ID,
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}

func mapContainerStatusesToApiModel(statuses []model.ContainerStatus) (result []*pb.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &pb.ContainerStatus{
			ContainerID: status.ContainerID,
			Image:       status.Image,
			State:       status.State,
		})
	}
	return result
}
