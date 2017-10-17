package mapping

import (
	log "github.com/sirupsen/logrus"

	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// MapModelByPodNamesToInternalModel builds from container label information pod data structure
func MapModelByPodNamesToInternalModel(containers []containerd.Container) map[string][]model.Container {
	result := make(map[string][]model.Container)
	for _, container := range containers {
		podName := GetPodName(container)

		if _, ok := result[podName]; !ok {
			result[podName] = []model.Container{}
		}

		result[podName] = append(result[podName], MapContainerToInternalModel(container))
	}
	return result
}

// GetPodName resolves pod name where the container belongs
func GetPodName(container containerd.Container) string {
	labels := ContainerLabels(container.Info().Labels)
	podName := labels.getPodName()
	if podName == "" {
		// container is not cand managed container so add it under 'system' pod in namespace 'default'
		podName = "system"
	}
	return podName
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
	spec, err := container.Spec()
	if err != nil {
		log.Fatalf("Cannot read container spec: %s", err)
		spec = &specs.Spec{}
	}
	return model.Container{
		Name:  container.ID(),
		Image: container.Info().Image,
		Tty:   spec.Process.Terminal,
	}
}

// MapContainerStatusToInternalModel maps containerd model to internal container status model
func MapContainerStatusToInternalModel(container containerd.Container, status containerd.Status) model.ContainerStatus {
	return model.ContainerStatus{
		ContainerID: container.ID(),
		Image:       container.Info().Image,
		State:       mapContainerStatus(status),
	}
}

func mapContainerStatus(status containerd.Status) string {
	if status.Status == "" {
		return string(containerd.Unknown)
	}
	return string(status.Status)
}
