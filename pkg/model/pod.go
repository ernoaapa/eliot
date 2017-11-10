package model

// DefaultNamespace is namespace what each pod get if there is no metadata.namespace
var DefaultNamespace = "elliot"

// Pod is set of containers
type Pod struct {
	Metadata Metadata `validate:"required"`
	Spec     PodSpec  `validate:"required"`
	Status   PodStatus
}

// PodSpec defines what containers should be running
type PodSpec struct {
	HostNetwork   bool
	HostPID       bool
	Containers    []Container `validate:"required,gt=0,dive"`
	RestartPolicy string
}

// PodStatus represents latest known state of pod
type PodStatus struct {
	Hostname          string
	ContainerStatuses []ContainerStatus `validate:"dive"`
}

// AppendContainer adds container to the pod information
func (p *Pod) AppendContainer(container Container, status ContainerStatus) {
	p.Spec.Containers = append(p.Spec.Containers, container)
	p.Status.ContainerStatuses = append(p.Status.ContainerStatuses, status)
}
