package model

// DeviceInfo contains information about current device
type DeviceInfo struct {
}

// DeviceState describes current state of the device
type DeviceState struct {
	Pods []PodState
}

// PodState represents information about pod current state
type PodState struct {
	ID string
}
