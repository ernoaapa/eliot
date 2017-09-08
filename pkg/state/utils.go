package state

import "github.com/ernoaapa/can/pkg/runtime"

func getPodUIDFromLabels(labels map[string]string) string {
	return runtime.GetLabelFor(labels, runtime.PodUIDSuffix)
}

func getPodNameFromLabels(labels map[string]string) string {
	return runtime.GetLabelFor(labels, runtime.PodNameSuffix)
}

func getPodNamespaceFromLabels(labels map[string]string) string {
	return runtime.GetLabelFor(labels, runtime.PodNamespaceSuffix)
}
func getContainerNameFromLabels(labels map[string]string) string {
	return runtime.GetLabelFor(labels, runtime.ContainerNameSuffix)
}
