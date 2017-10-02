package model

// DefaultNamespace is namespace what each pod get if there is no metadata.namespace
var DefaultNamespace = "cand"

// Pod is set of containers
type Pod struct {
	Metadata Metadata `validate:"required"`
	Spec     PodSpec  `validate:"required"`
	Status   PodStatus
}

// PodSpec defines what containers should be running
type PodSpec struct {
	HostNetwork bool
	Containers  []Container `validate:"required,gt=0,dive"`
}

// PodStatus represents latest known state of pod
type PodStatus struct {
	ContainerStatuses []ContainerStatus `validate:"dive"`
}
