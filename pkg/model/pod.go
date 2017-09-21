package model

// Pod is set of containers
type Pod struct {
	Metadata
	UID    string    `                    json:"uid"     yaml:"uid"`
	Spec   PodSpec   `validate:"required" json:"spec"     yaml:"spec"`
	Status PodStatus `                    json:"status"   yaml:"status"`
}

// PodSpec defines what containers should be running
type PodSpec struct {
	Containers []Container `validate:"required,gt=0,dive"  json:"containers"  yaml:"containers"`
}

// PodStatus represents latest known state of pod
type PodStatus struct {
	ContainerStatuses []ContainerStatus `validate:"dive"   json:"containerStatuses,omitempty"   yaml:"containerStatuses,omitempty"`
}
