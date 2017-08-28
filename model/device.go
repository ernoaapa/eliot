package model

// DeviceInfo contains information about current device
type DeviceInfo struct {
	// Labels for the device, provided through cli
	Labels map[string]string `json:"labels"`

	// The machine id is an ID identifying a specific Linux/Unix installation.
	// It does not change if hardware is replaced.
	MachineID string `json:"machine_id"`

	// The system uuid is the main board product UUID,
	// as set by the board manufacturer and encoded in the BIOS DMI information
	SystemUUID string `json:"system_uuid"`

	// A random ID that is regenerated on each boot
	BootID string `json:"boot_id"`

	// Device operating system. One of 386, amd64, arm, s390x, and so on.
	Arch string `json:"arch"`

	// device operating system. One of darwin, freebsd, linux, and so on
	OS string `json:"os"`

	// Device kernel e.g. "Linux"
	Kernel string `json:"kernel"`

	// Device core version e.g. "3.13.0-27-generic"
	Core string `json:"core"`

	// Device platform e.g. x86_64
	Platform string `json:"platform"`

	// Device hostname
	Hostname string `json:"hostname"`

	// Number of CPUs
	CPUs int `json:"cpus"`
}

// DeviceState describes current state of the device
type DeviceState struct {
	Pods []PodState
}

// PodState represents information about pod current state
type PodState struct {
	ID string
}
