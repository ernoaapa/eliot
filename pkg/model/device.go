package model

import (
	"fmt"
	"net"
)

// DeviceInfo contains information about current device
type DeviceInfo struct {
	// Labels for the device, provided through cli
	Labels map[string]string

	// Device hostname
	Hostname string `validate:"required,gt=0"`

	// IPs
	Addresses []net.IP

	// Port
	GrpcPort int

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

// GetPrimaryEndpoint return primary GRPC endpoint address
func (d DeviceInfo) GetPrimaryEndpoint() string {
	if len(d.Addresses) == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", d.Addresses[0], d.GrpcPort)
}
