package printers

import (
	"bytes"
	"testing"

	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/stretchr/testify/assert"
)

var examplePod = &pods.Pod{
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
}

func TestSuite(t *testing.T) {
	implementations := map[string]ResourcePrinter{
		"human": NewHumanReadablePrinter(),
		"yaml":  NewYamlPrinter(),
	}

	for name, impl := range implementations {
		t.Run(name, func(t *testing.T) {
			testYamlPrintPodsTable(t, impl)
			testYamlPrintDevicesTable(t, impl)
			testPrintDeviceDetails(t, impl)
			testPrintPodsTable(t, impl)
			testPrintConfig(t, impl)
		})
	}
}

func testYamlPrintPodsTable(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	err := printer.PrintPodsTable([]*pods.Pod{examplePod}, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testYamlPrintDevicesTable(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	data := []model.DeviceInfo{
		{
			Hostname: "foobar",
			Labels: map[string]string{
				"env": "test",
			},
		},
	}

	err := printer.PrintDevicesTable(data, &buffer)
	assert.NoError(t, err, "Printing devices table should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testPrintDeviceDetails(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

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

func testPrintPodsTable(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	err := printer.PrintPodsTable([]*pods.Pod{examplePod}, &buffer)
	assert.NoError(t, err, "Printing pods table should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testPrintPodDetails(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

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
func testPrintConfig(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

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
