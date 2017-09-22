package containerd

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

// ContainerLabels is helper type for managing container labels
type ContainerLabels map[string]string

func (l ContainerLabels) getPodName() string {
	return l.getValue(PodNameSuffix)
}

func (l ContainerLabels) getContainerName() string {
	return l.getValue(ContainerNameSuffix)
}

func (l ContainerLabels) getValue(key string) string {
	return l[getLabelKeyFor(key)]
}

func getLabelKeyFor(name string) string {
	return fmt.Sprintf("%s.%s", ContainerLabelPrefix, name)
}

// NewContainerLabels constructs new labels map for new container
func NewContainerLabels(pod model.Pod, container model.Container) ContainerLabels {
	labels := make(map[string]string)
	labels[getLabelKeyFor(PodUIDSuffix)] = pod.UID
	labels[getLabelKeyFor(PodNameSuffix)] = pod.GetName()
	labels[getLabelKeyFor(PodNamespaceSuffix)] = pod.GetNamespace()
	labels[getLabelKeyFor(ContainerNameSuffix)] = container.Name
	return labels
}
