package containerd

import (
	"fmt"

	"github.com/ernoaapa/can/pkg/model"
)

var (
	// LabelPrefix is prefix what all container labels what cand creates get
	labelPrefix        = "io.can"
	podNameLabel       = "pod.name"
	containerNameLabel = "container.name"
)

// ContainerLabels is helper type for managing container labels
type ContainerLabels map[string]string

func (l ContainerLabels) getPodName() string {
	return l.getValue(podNameLabel)
}

func (l ContainerLabels) getContainerName() string {
	return l.getValue(containerNameLabel)
}

func (l ContainerLabels) getValue(key string) string {
	return l[buildLabelKeyFor(key)]
}

func buildLabelKeyFor(name string) string {
	return fmt.Sprintf("%s.%s", labelPrefix, name)
}

// NewLabels constructs new labels map for new container
func NewLabels(pod model.Pod, container model.Container) ContainerLabels {
	labels := make(map[string]string)
	labels[buildLabelKeyFor(podNameLabel)] = pod.Metadata.Name
	labels[buildLabelKeyFor(containerNameLabel)] = container.Name
	return labels
}
