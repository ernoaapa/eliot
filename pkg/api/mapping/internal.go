package mapping

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

func MapPodsToAPIModel(namespace string, containersByPods map[string][]model.Container) (result []*pb.Pod) {
	for podName, containers := range containersByPods {
		result = append(result, MapPodToAPIModel(namespace, podName, containers))
	}
	return result
}

func MapPodToAPIModel(namespace, podName string, containers []model.Container) *pb.Pod {
	return &pb.Pod{
		Metadata: model.NewMetadata(podName, namespace),
		Spec: &pb.PodSpec{
			Containers: MapContainersToAPIModel(containers),
		},
		// Status: &pb.PodStatus{
		// 	ContainerStatuses: mapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
		// },
	}
}

func MapContainersToAPIModel(containers []model.Container) (result []*pb.Container) {
	for _, container := range containers {
		result = append(result, &pb.Container{
			ID:    container.ID,
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}

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
