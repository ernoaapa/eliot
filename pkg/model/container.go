package model

// Container defines what image should be running
type Container struct {
	Name  string `validate:"required,gt=0,alphanumOrDash"   json:"name"      yaml:"name"`
	Image string `validate:"required,gt=0,imageRef"         json:"image"     yaml:"image"`
}

// ContainerStatus represents one container status
type ContainerStatus struct {
	ContainerID string `json:"containerId,omitempty"  yaml:"containerId,omitempty"`
	Image       string `json:"image,omitempty"        yaml:"image,omitempty"`
	State       string `json:"state,omitempty"        yaml:"state,omitempty"`
}
