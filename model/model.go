package model

import "fmt"

var (
	nameFieldName = "name"
)

// Pod is set of containers
type Pod struct {
	Metadata map[string]string `yaml:"metadata"`
	Spec     Spec              `yaml:"spec"`
}

// GetName returns name from metadata
func (p *Pod) GetName() string {
	return p.Metadata[nameFieldName]
}

// Spec defines what containers should be running
type Spec struct {
	Containers []Container `yaml:"containers"`
}

// Container defines what image should be running
type Container struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

// BuildID creates unique id for the container from parent pod name
func (c *Container) BuildID(podName string) string {
	return fmt.Sprintf("%s-%s", podName, c.Name)
}

// DeviceInfo contains information about current device
type DeviceInfo struct {
}
