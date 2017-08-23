package model

import (
	"fmt"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Pod is set of containers
type Pod struct {
	Metadata Metadata `validate:"hasName" yaml:"metadata"`
	Spec     Spec     `validate:"required" yaml:"spec"`
}

// GetName returns pod name from metadata
func (p *Pod) GetName() string {
	return p.Metadata.GetName()
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

// Spec defines what containers should be running
type Spec struct {
	Containers []Container `validate:"required,gt=0,dive" yaml:"containers"`
}

// Container defines what image should be running
type Container struct {
	Name  string `validate:"required,gt=0,alphanum" yaml:"name"`
	Image string `validate:"required,gt=0,imageRef" yaml:"image"`
}

// BuildID creates unique id for the container from parent pod name
func (c *Container) BuildID(podName string) string {
	return fmt.Sprintf("%s-%s", podName, c.Name)
}
