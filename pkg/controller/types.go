package controller

import (
	"github.com/ernoaapa/can/pkg/model"
)

type podsManifest []model.Pod

func (list podsManifest) containsPod(name string) bool {
	for _, pod := range list {
		if pod.Metadata.Name == name {
			return true
		}
	}
	return false
}

func (list podsManifest) containsContainer(podName string, target model.Container) bool {
	for _, pod := range list {
		if pod.Metadata.Name == podName {
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
		if pod.Metadata.Namespace == namespace {
			result = append(result, pod)
		}
	}
	return result
}

func (list podsManifest) getNamespaces() []string {
	result := []string{}
	for _, pod := range list {
		namespace := pod.Metadata.Namespace
		if namespace != "" {
			result = append(result, pod.Metadata.Namespace)
		}
	}
	return result
}

type podsState []model.Pod

func (p podsState) getPodContainers(podName string) []model.Container {
	for _, pod := range p {
		if pod.Metadata.Name == podName {
			return pod.Spec.Containers
		}
	}
	return []model.Container{}
}

func (p podsState) containsContainer(podName string, target model.Container) bool {
	for _, container := range p.getPodContainers(podName) {
		if containersMatch(container, target) {
			return true
		}
	}
	return false
}

func (p podsState) findContainer(podName string, target model.Container) string {
	for _, container := range p.getPodContainers(podName) {
		if containersMatch(container, target) {
			return container.Name
		}
	}
	return ""
}

func containersMatch(container model.Container, target model.Container) bool {
	return container.Name == target.Name &&
		container.Image == target.Image
}
