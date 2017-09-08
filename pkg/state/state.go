package state

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
)

func getCurrentState(client *runtime.ContainerdClient) (result []*model.Pod, err error) {
	result = []*model.Pod{}
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return result, err
	}

	for _, namespace := range namespaces {
		containers, err := client.GetContainers(namespace)
		if err != nil {
			return result, err
		}

		result = append(result, constructPodsFromContainerInfo(client, containers)...)
	}
	return result, nil
}

func constructPodsFromContainerInfo(client *runtime.ContainerdClient, containers []containerd.Container) []*model.Pod {
	podsByName := make(map[string]*model.Pod)

	for _, container := range containers {
		labels := container.Info().Labels
		podName := getPodNameFromLabels(labels)
		podNamespace := getPodNamespaceFromLabels(labels)
		if _, ok := podsByName[podName]; !ok {
			podsByName[podName] = &model.Pod{
				UID: getPodUIDFromLabels(labels),
				Metadata: model.Metadata{
					"name":      podName,
					"namespace": podNamespace,
				},
				Spec: model.Spec{
					Containers: []model.Container{},
				},
				Status: model.PodStatus{
					ContainerStatuses: []model.ContainerStatus{},
				},
			}
		}
		podsByName[podName].Spec.Containers = append(podsByName[podName].Spec.Containers, model.Container{
			Name:  getContainerNameFromLabels(labels),
			Image: container.Info().Image,
		})

		podsByName[podName].Status.ContainerStatuses = append(podsByName[podName].Status.ContainerStatuses, resolveContainerStatus(client, container))
	}
	return getValues(podsByName)
}

func resolveContainerStatus(client *runtime.ContainerdClient, container containerd.Container) model.ContainerStatus {
	return model.ContainerStatus{
		ContainerID: container.ID(),
		Image:       container.Info().Image,
		State:       client.GetContainerTaskStatus(container.ID()),
	}
}

func getValues(podsByName map[string]*model.Pod) []*model.Pod {
	values := []*model.Pod{}
	for _, pod := range podsByName {
		values = append(values, pod)
	}
	return values
}
