package pods

import (
	containers "github.com/ernoaapa/elliot/pkg/api/services/containers/v1"
	"github.com/ernoaapa/elliot/pkg/model"
)

// Defaults set default values to pod definitions
func Defaults(pods []*Pod) (result []*Pod) {
	for _, pod := range pods {
		result = append(result, Default(pod))
	}
	return result
}

// Default set default values to Pod model
func Default(pod *Pod) *Pod {
	if pod.Metadata.Namespace == "" {
		pod.Metadata.Namespace = model.DefaultNamespace
	}

	pod.Spec.Containers = containers.Defaults(pod.Spec.Containers)
	return pod
}
