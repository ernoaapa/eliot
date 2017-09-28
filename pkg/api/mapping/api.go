package mapping

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

// MapPodsToInternalModel maps API Pod model to internal model
func MapPodsToInternalModel(pods []*pb.Pod) (result []model.Pod) {
	for _, pod := range pods {
		result = append(result, MapPodToInternalModel(pod))
	}
	return result
}

// MapPodToInternalModel maps API Pod model to internal model
func MapPodToInternalModel(pod *pb.Pod) model.Pod {
	return model.Pod{
		Metadata: model.Metadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
		},
		Spec: model.PodSpec{
			Containers: MapContainerToInternalModel(pod.Spec.Containers),
		},
	}
}

// MapContainerToInternalModel maps API Container model to internal model
func MapContainerToInternalModel(containers []*pb.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, model.Container{
			Name:  container.Name,
			Image: container.Image,
		})
	}
	return result
}
