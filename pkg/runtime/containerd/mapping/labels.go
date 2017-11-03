package mapping

import (
	"fmt"

	"github.com/ernoaapa/can/pkg/model"
)

var (
	// LabelPrefix is prefix what all container labels what cand creates get
	labelPrefix  = "io.can"
	podNameLabel = "pod.name"
	stdinLabel   = "pod.io.stdin"
	stdoutLabel  = "pod.io.stdout"
	stderrLabel  = "pod.io.stderr"
)

// ContainerLabels is helper type for managing container labels
type ContainerLabels map[string]string

func (l ContainerLabels) getPodName() string {
	return l.getValue(podNameLabel)
}

func (l ContainerLabels) getIoSet() model.IOSet {
	return model.IOSet{
		In:  l.getValue(stdinLabel),
		Out: l.getValue(stdoutLabel),
		Err: l.getValue(stderrLabel),
	}
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
	labels[buildLabelKeyFor(stdinLabel)] = container.Io.In
	labels[buildLabelKeyFor(stdoutLabel)] = container.Io.Out
	labels[buildLabelKeyFor(stderrLabel)] = container.Io.Err
	return labels
}
