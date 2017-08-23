package model

// Pod is set of containers
type Pod struct {
	Name string `yaml:"name"`
	Spec Spec   `yaml:"spec"`
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

// DeviceInfo contains information about current device
type DeviceInfo struct {
}
