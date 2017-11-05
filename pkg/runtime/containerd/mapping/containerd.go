package mapping

import (
	log "github.com/sirupsen/logrus"

	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime/containerd/extensions"
)

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
	labels := ContainerLabels(container.Info().Labels)
	return model.Container{
		Name:  labels.getContainerName(),
		Image: container.Info().Image,
		Tty:   RequireTty(container),
		Pipe:  mapPipeToInternalModel(container),
	}
}

// RequireTty find out is the container configured to create TTY
func RequireTty(container containerd.Container) bool {
	spec, err := container.Spec()
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process TTY value: %s", err)
		return false
	}
	return spec.Process.Terminal
}

func mapPipeToInternalModel(container containerd.Container) *model.PipeSet {
	pipe, err := extensions.GetPipeExtension(container)
	if err != nil {
		log.Errorf("Failed to read Pipe extension from container [%s]: %s", container.ID(), err)
	}
	if pipe == nil {
		return nil
	}

	return &model.PipeSet{
		Stdout: &model.PipeFromStdout{
			Stdin: &model.PipeToStdin{
				Name: pipe.Stdout.Stdin.Name,
			},
		},
	}
}

// MapContainerStatusToInternalModel maps containerd model to internal container status model
func MapContainerStatusToInternalModel(container containerd.Container, status containerd.Status) model.ContainerStatus {
	labels := ContainerLabels(container.Info().Labels)
	return model.ContainerStatus{
		ContainerID: container.ID(),
		Name:        labels.getContainerName(),
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
