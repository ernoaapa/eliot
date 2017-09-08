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
	return labels[getLabelKeyFor(key)]
}

func getLabelKeyFor(name string) string {
	return fmt.Sprintf("%s.%s", ContainerLabelPrefix, name)
}

func getContainerLabels(pod model.Pod, container model.Container) map[string]string {
	labels := make(map[string]string)
	labels[getLabelKeyFor(PodUIDSuffix)] = pod.UID
	labels[getLabelKeyFor(PodNameSuffix)] = pod.GetName()
	labels[getLabelKeyFor(PodNamespaceSuffix)] = pod.GetNamespace()
	labels[getLabelKeyFor(ContainerNameSuffix)] = container.Name
	return labels
}
