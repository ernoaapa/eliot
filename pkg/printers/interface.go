package printers

import (
	"io"

	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/model"
)

// ResourcePrinter is an interface that knows how to print runtime objects.
type ResourcePrinter interface {
	// PrintPodsTable receives a list of Pods, formats it as table and prints it to a writer.
	PrintPodsTable([]*pods.Pod, io.Writer) error
	// PrintDevicesTable receives a channel of Devices, formats the output as table
	PrintDevicesTable([]model.DeviceInfo, io.Writer) error
	// PrintDeviceDetails receives a Device, formats it to detailed description and prints it to the writer.
	PrintDeviceDetails(*device.Info, io.Writer) error
	// PrintPodDetails receives a Pod, formats it to detailed description and prints it to the writer.
	PrintPodDetails(*pods.Pod, io.Writer) error
	// PrintConfig receives config and formats it to human readable format
	PrintConfig(*config.Config, io.Writer) error
}
