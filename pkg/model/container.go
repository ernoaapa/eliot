package model

// Container defines what image should be running
type Container struct {
	Name       string `validate:"required,gt=0,alphanumOrDash"`
	Image      string `validate:"required,gt=0,imageRef"`
	Tty        bool
	Args       []string `validate:"dive,noSpaces"`
	Env        []string `validate:"dive,envKeyValuePair"`
	WorkingDir string   `validate:"omitempty,gt=0"`
	Mounts     []Mount  `validate:"dive"`
	Io         IOSet
}

type IOSet struct {
	In  string
	Out string
	Err string
}

// Mount defines directory mount from host to the container
type Mount struct {
	Type        string   `validate:"omitempty,gt=0"`
	Source      string   `validate:"omitempty,gt=0"`
	Destination string   `validate:"omitempty,gt=0"`
	Options     []string `validate:"dive,gt=0"`
}

// ContainerStatus represents one container status
type ContainerStatus struct {
	ContainerID string `validate:"required,gt=0"`
	Image       string `validate:"required,gt=0,imageRef"`
	State       string `validate:"required,gt=0"`
}
