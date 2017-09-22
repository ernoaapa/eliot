package controller

import (
	"github.com/ernoaapa/can/pkg/model"
)

type podsManifest []model.Pod

func (list podsManifest) containsPod(name string) bool {
	for _, pod := range list {
		if pod.GetName() == name {
			return true
		}
	}
	return false
}

func (list podsManifest) containsContainer(podName string, target model.Container) bool {
	for _, pod := range list {
		if pod.GetName() == podName {
			for _, container := range pod.Spec.Containers {
				if containersMatch(container, target) {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (list podsManifest) filterPodsByNamespace(namespace string) podsManifest {
	result := podsManifest{}
	for _, pod := range list {
		if pod.GetNamespace() == namespace {
			result = append(result, pod)
		}
	}
	return result
}

func (list podsManifest) getNamespaces() []string {
	result := []string{}
	for _, pod := range list {
		namespace := pod.GetNamespace()
		if namespace != "" {
			result = append(result, pod.GetNamespace())
		}
	}
	return result
}

type containersState map[string][]model.Container

func (c containersState) getPodContainers(podName string) []model.Container {
	return c[podName]
}

func (c containersState) containsContainer(podName string, target model.Container) bool {
	for _, container := range c.getPodContainers(podName) {
		if containersMatch(container, target) {
			return true
		}
	}
	return false
}

func containersMatch(container model.Container, target model.Container) bool {
	return container.Name == target.Name &&
		container.Image == target.Image
}
