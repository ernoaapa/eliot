package model

import (
	"net"
)

// NodeInfo contains information about current node
type NodeInfo struct {
	// Labels for the node, provided through cli
	Labels map[string]string

	// Node hostname
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

	// Node operating system. One of 386, amd64, arm, s390x, and so on.
	Arch string

	// node operating system. One of darwin, freebsd, linux, windows, and so on
	OS string

	// Server version
	Version string

	// Filesystems
	Filesystems []Filesystem

	// Seconds since node boot up
	Uptime uint64
}

// NodeState describes current state of the node
type NodeState struct {
	Pods []PodState `validate:"dive"`
}

// PodState represents information about pod current state
type PodState struct {
	ID string `validate:"required,gt=0"`
}

// Filesystem represents information about single filesystem in the target node
type Filesystem struct {
	// E.g /dev/vda1, tmpfs, cgroup, etc.
	Filesystem string
	// E.g. ext4, tmpfs, cgroup, etc.
	TypeName string
	// Path to the mounted dir
	MountDir string
	// Total data blocks
	Total uint64
	// Free blocks
	Free uint64
	// Free blocks available to unprivileged user
	Available uint64
}
