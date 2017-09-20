package model

import (
	"sync"

	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Pod is set of containers
type Pod struct {
	UID      string    `                    json:"uid"     yaml:"uid"`
	Metadata Metadata  `validate:"hasName"  json:"metadata" yaml:"metadata"`
	Spec     Spec      `validate:"required" json:"spec"     yaml:"spec"`
	Status   PodStatus `                    json:"status"   yaml:"status"`
}

// GetName returns pod name from metadata
func (p *Pod) GetName() string {
	return p.Metadata.GetName()
}

// GetNamespace returns pod namespace from metadata
func (p *Pod) GetNamespace() string {
	return p.Metadata.GetNamespace()
}

// Metadata is any random metadata
type Metadata map[string]interface{}

// GetName returns name field from metadata
func (m Metadata) GetName() string {
	name := m["name"]
	if name == nil {
		return ""
	}
	return m["name"].(string)
}

// GetNamespace returns name field from metadata
func (m Metadata) GetNamespace() string {
	namespace := m["namespace"]
	if namespace == nil {
		return ""
	}
	return m["namespace"].(string)
}

// Spec defines what containers should be running
type Spec struct {
	Containers []Container `validate:"required,gt=0,dive"  json:"containers"  yaml:"containers"`
}

// Container defines what image should be running
type Container struct {
	ID    string
	Name  string `validate:"required,gt=0,alphanumOrDash"   json:"name"      yaml:"name"`
	Image string `validate:"required,gt=0,imageRef"         json:"image"     yaml:"image"`
}

// BuildContainerID creates unique id for the container
func BuildContainerID() string {
	return uuid.NewV4().String()
}

// PodStatus represents latest known state of pod
type PodStatus struct {
	ContainerStatuses []ContainerStatus `validate:"dive"   json:"containerStatuses,omitempty"   yaml:"containerStatuses,omitempty"`
}

// ContainerStatus represents one container status
type ContainerStatus struct {
	ContainerID string `json:"containerId,omitempty"  yaml:"containerId,omitempty"`
	Image       string `json:"image,omitempty"        yaml:"image,omitempty"`
	State       string `json:"state,omitempty"        yaml:"state,omitempty"`
}
