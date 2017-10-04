package mapping

import (
	"log"

	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// MapModelByPodNamesToInternalModel builds from container label information pod data structure
func MapModelByPodNamesToInternalModel(containers []containerd.Container) map[string][]model.Container {
	result := make(map[string][]model.Container)
	for _, container := range containers {
		labels := ContainerLabels(container.Info().Labels)
		podName := labels.getPodName()
		if podName == "" {
			// container is not cand managed container so add it under 'system' pod in namespace 'default'
			podName = "system"
		}

		if _, ok := result[podName]; !ok {
			result[podName] = []model.Container{}
		}

		result[podName] = append(result[podName], MapContainerToInternalModel(container))
	}
	return result
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
