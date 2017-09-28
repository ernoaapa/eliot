package model

// DeviceInfo contains information about current device
type DeviceInfo struct {
	// Labels for the device, provided through cli
	Labels map[string]string

	// The machine id is an ID identifying a specific Linux/Unix installation.
	// It does not change if hardware is replaced.
	MachineID string `validate:"required,gt=0"`

	// The system uuid is the main board product UUID,
	// as set by the board manufacturer and encoded in the BIOS DMI information
	SystemUUID string `validate:"required,gt=0"`

	// A random ID that is regenerated on each boot
	BootID string `validate:"required,gt=0"`

	// Device operating system. One of 386, amd64, arm, s390x, and so on.
	Arch string

	// device operating system. One of darwin, freebsd, linux, and so on
	OS string

	// Device kernel e.g. "Linux"
	Kernel string

	// Device core version e.g. "3.13.0-27-generic"
	Core string

	// Device platform e.g. x86_64
	Platform string

	// Device hostname
	Hostname string `validate:"required,gt=0"`

	// Number of CPUs
	CPUs int
}

// DeviceState describes current state of the device
type DeviceState struct {
	Pods []PodState `validate:"dive"`
}

// PodState represents information about pod current state
type PodState struct {
	ID string `validate:"required,gt=0"`
}
