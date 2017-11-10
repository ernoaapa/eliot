package api

import (
	"fmt"

	containers "github.com/ernoaapa/elliot/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/elliot/pkg/api/services/pods/v1"
)

// WithSharedMount adds mount point to each container
func WithSharedMount(mount *containers.Mount) PodOpts {
	return func(pod *pods.Pod) error {
		for _, container := range pod.Spec.Containers {
			container.Mounts = append(container.Mounts, mount)
		}
		return nil
	}
}

// WithWorkingDir sets each container workdir
// If workdir is already defined, will return error
func WithWorkingDir(workDir string) PodOpts {
	return func(pod *pods.Pod) error {
		for _, container := range pod.Spec.Containers {
			if container.WorkingDir != "" {
				return fmt.Errorf("Container [%s] already have WorkingDir defined", container.Name)
			}
			container.WorkingDir = workDir
		}
		return nil
	}
}

// WithContainer adds container to the Pod spec
func WithContainer(container *containers.Container) PodOpts {
	return func(pod *pods.Pod) error {
		pod.Spec.Containers = append(pod.Spec.Containers, container)
		return nil
	}
}
