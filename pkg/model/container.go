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
	Devices	   []Device `validate:"dive"`
	Pipe       *PipeSet
}

// PipeSet allows defining pipe from some source(s) to another container
type PipeSet struct {
	Stdout *PipeFromStdout
}

// PipeFromStdout defines container stdout as source for the piping
type PipeFromStdout struct {
	Stdin *PipeToStdin
}

// PipeToStdin defines container stdin as target for the piping
type PipeToStdin struct {
	Name string
}

// Mount defines directory mount from host to the container
type Mount struct {
	Type        string   `validate:"omitempty,gt=0"`
	Source      string   `validate:"omitempty,gt=0"`
	Destination string   `validate:"omitempty,gt=0"`
	Options     []string `validate:"dive,gt=0"`
}

type Device struct {
	DeviceType	string
	MajorId		uint32
	MinorId		uint32
}

// ContainerStatus represents one container status
type ContainerStatus struct {
	ContainerID  string `validate:"required,gt=0"`
	Name         string `validate:"required,gt=0"`
	Image        string `validate:"required,gt=0,imageRef"`
	State        string `validate:"required,gt=0"`
	RestartCount int    `validate:"required,gte=0"`
}
