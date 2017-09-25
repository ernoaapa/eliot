package model

import uuid "github.com/satori/go.uuid"

// DefaultNamespace is namespace what each pod get if there is no metadata.namespace
var DefaultNamespace = "cand"

// Defaults set default values to pod definitions
func Defaults(pods []Pod) (result []Pod) {
	for _, pod := range pods {
		result = append(result, Default(pod))
	}
	return result
}

// Default set default values to Pod model
func Default(pod Pod) Pod {
	if pod.GetNamespace() == "" {
		pod.SetNamespace(DefaultNamespace)
	}

	pod.Spec.Containers = defaultContainers(pod.GetName(), pod.Spec.Containers)
	return pod
}

func defaultContainers(podName string, containers []Container) (result []Container) {
	for _, container := range containers {
		container.ID = BuildContainerID()
		result = append(result, container)
	}
	return result
}

// BuildContainerID creates unique id for the container
func BuildContainerID() string {
	return uuid.NewV4().String()
}
