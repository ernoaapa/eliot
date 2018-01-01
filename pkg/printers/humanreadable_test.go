package printers

import (
	"bytes"
	"testing"

	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPrintDeviceDetails(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := &device.Info{
		Labels:     []*device.Label{&device.Label{Key: "foo", Value: "bar"}},
		Hostname:   "foo-bar",
		Addresses:  []string{"1.2.3.4"},
		GrpcPort:   5000,
		MachineID:  "1234-5678",
		SystemUUID: "asdf-jklö",
		BootID:     "12334345345",
		Arch:       "amd64",
		Os:         "linux",
	}

	err := printer.PrintDeviceDetails(data, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func TestPrintTable(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := []*pods.Pod{
		&pods.Pod{
			Metadata: &core.ResourceMetadata{
				Name:      "foo",
				Namespace: "eliot",
			},
			Spec: &pods.PodSpec{
				Containers: []*containers.Container{
					&containers.Container{},
					&containers.Container{},
				},
			},
		},
	}

	err := printer.PrintPodsTable(data, &buffer)
	assert.NoError(t, err, "Printing pods table should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func TestPrintPodDetails(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      "foo",
			Namespace: "eliot",
		},
		Spec: &pods.PodSpec{
			Containers: []*containers.Container{
				&containers.Container{},
				&containers.Container{},
			},
		},
		Status: &pods.PodStatus{
			Hostname: "testing.local",
		},
	}

	err := printer.PrintPodDetails(data, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
func TestPrintConfig(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := &config.Config{
		Endpoints: []config.Endpoint{
			config.Endpoint{Name: "localhost", URL: "localhost:5000"},
		},
		Namespace: "default",
	}

	err := printer.PrintConfig(data, &buffer)
	assert.NoError(t, err, "Printing config should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
