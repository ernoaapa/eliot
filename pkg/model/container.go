package model

// Container defines what image should be running
type Container struct {
	Name  string `validate:"required,gt=0,alphanumOrDash"`
	Image string `validate:"required,gt=0,imageRef"`
}

// ContainerStatus represents one container status
type ContainerStatus struct {
	ContainerID string `validate:"required,gt=0"`
	Image       string `validate:"required,gt=0,imageRef"`
	State       string `validate:"required,gt=0"`
}
