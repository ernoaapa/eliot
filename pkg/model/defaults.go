package model

// DefaultNamespace is namespace what each pod get if there is no metadata.namespace
var DefaultNamespace = "cand"

// Defaults set default values to pod definitions
func Defaults(pods []Pod) (result []Pod) {
	for _, pod := range pods {
		result = append(result, Default(pod))
	}
	return result
}

// Default set default values to Pod model
func Default(pod Pod) Pod {
	if pod.Metadata.Namespace == "" {
		pod.Metadata.Namespace = DefaultNamespace
	}
	return pod
}
