package runtime

import (
	"fmt"

	"github.com/ernoaapa/can/pkg/model"
)

var (
	ContainerLabelPrefix = "io.can"
	PodUIDSuffix         = "pod.uid"
	PodNameSuffix        = "pod.name"
	PodNamespaceSuffix   = "pod.namespace"
	ContainerNameSuffix  = "container.name"
)

func GetLabelFor(labels map[string]string, key string) string {
	return labels[GetLabelKeyFor(key)]
}

func GetLabelKeyFor(name string) string {
	return fmt.Sprintf("%s.%s", ContainerLabelPrefix, name)
}

func getContainerLabels(pod model.Pod, container model.Container) map[string]string {
	labels := make(map[string]string)
	labels[GetLabelKeyFor(PodUIDSuffix)] = pod.UID
	labels[GetLabelKeyFor(PodNameSuffix)] = pod.GetName()
	labels[GetLabelKeyFor(PodNamespaceSuffix)] = pod.GetNamespace()
	labels[GetLabelKeyFor(ContainerNameSuffix)] = container.Name
	return labels
}
