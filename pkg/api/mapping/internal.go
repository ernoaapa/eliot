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

// MapPodToAPIModel maps internal Pod model to API model
func MapPodToAPIModel(pod model.Pod) *pods.Pod {
	return &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
		},
		Spec: &pods.PodSpec{
			Containers:  MapContainersToAPIModel(pod.Spec.Containers),
			HostNetwork: pod.Spec.HostNetwork,
			HostPID:     pod.Spec.HostPID,
		},
		Status: &pods.PodStatus{
			Hostname:          pod.Status.Hostname,
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
			Pipe:  mapPipeToAPIModel(container.Pipe),
		})
	}
	return result
}

func mapPipeToAPIModel(pipe *model.PipeSet) *containers.PipeSet {
	if pipe == nil {
		return nil
	}
	return &containers.PipeSet{
		Stdout: &containers.PipeFromStdout{
			Stdin: &containers.PipeToStdin{
				Name: pipe.Stdout.Stdin.Name,
			},
		},
	}
}

// MapContainerStatusesToAPIModel maps list of internal ContainerStatus models to API model
func MapContainerStatusesToAPIModel(statuses []model.ContainerStatus) (result []*containers.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &containers.ContainerStatus{
			ContainerID:  status.ContainerID,
			Name:         status.Name,
			Image:        status.Image,
			State:        status.State,
			RestartCount: int32(status.RestartCount),
		})
	}
	return result
}
