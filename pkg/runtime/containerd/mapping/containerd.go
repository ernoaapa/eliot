package mapping

import (
	"encoding/json"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime/containerd/extensions"
)

// GetPodName resolves pod name where the container belongs
func GetPodName(container containers.Container) string {
	labels := ContainerLabels(container.Labels)
	podName := labels.getPodName()
	if podName == "" {
		// container is not cand managed container so add it under 'system' pod in namespace 'default'
		podName = "system"
	}
	return podName
}

// MapContainersToInternalModel maps containerd models to internal model
func MapContainersToInternalModel(containers []containers.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, MapContainerToInternalModel(container))
	}
	return result
}

// MapContainerToInternalModel maps containerd model to internal model
func MapContainerToInternalModel(container containers.Container) model.Container {
	labels := ContainerLabels(container.Labels)
	return model.Container{
		Name:  labels.getContainerName(),
		Image: container.Image,
		Tty:   RequireTty(container),
		Pipe:  mapPipeToInternalModel(container),
	}
}

// RequireTty find out is the container configured to create TTY
func RequireTty(container containers.Container) bool {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process TTY value: %s", err)
		return false
	}
	return spec.Process.Terminal
}

// Spec returns the current OCI specification for the container
func getSpec(container containers.Container) (*specs.Spec, error) {
	var s specs.Spec
	if err := json.Unmarshal(container.Spec.Value, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func mapPipeToInternalModel(container containers.Container) *model.PipeSet {
	pipe, err := extensions.GetPipeExtension(container)
	if err != nil {
		log.Errorf("Failed to read Pipe extension from container [%s]: %s", container.ID, err)
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
func MapContainerStatusToInternalModel(container containers.Container, status containerd.Status) model.ContainerStatus {
	labels := ContainerLabels(container.Labels)
	return model.ContainerStatus{
		ContainerID: container.ID,
		Name:        labels.getContainerName(),
		Image:       container.Image,
		State:       mapContainerStatus(status),
	}
}

func mapContainerStatus(status containerd.Status) string {
	if status.Status == "" {
		return string(containerd.Unknown)
	}
	return string(status.Status)
}
