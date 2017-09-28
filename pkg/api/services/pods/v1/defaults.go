package pods

import "github.com/ernoaapa/can/pkg/model"

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
	return pod
}
