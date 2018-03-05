package printers

import (
	"bytes"
	"testing"

	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/stretchr/testify/assert"
)

var examplePod = &pods.Pod{
	Metadata: &core.ResourceMetadata{
		Name:      "foo",
		Namespace: "eliot",
	},
	Spec: &pods.PodSpec{
		Containers: []*containers.Container{
			{},
			{},
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
			testYamlPrintPods(t, impl)
			testYamlPrintNodes(t, impl)
			testPrintNode(t, impl)
			testPrintPods(t, impl)
			testPrintConfig(t, impl)
		})
	}
}

func testYamlPrintPods(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	err := printer.PrintPods([]*pods.Pod{examplePod}, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testYamlPrintNodes(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	data := []*node.Info{
		{
			Hostname: "foobar",
			Labels: []*node.Label{
				{Key: "env", Value: "test"},
			},
		},
	}

	err := printer.PrintNodes(data, &buffer)
	assert.NoError(t, err, "Printing nodes table should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testPrintNode(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	data := &node.Info{
		Labels:     []*node.Label{{Key: "foo", Value: "bar"}},
		Hostname:   "foo-bar",
		Addresses:  []string{"1.2.3.4"},
		GrpcPort:   5000,
		MachineID:  "1234-5678",
		SystemUUID: "asdf-jklö",
		BootID:     "12334345345",
		Arch:       "amd64",
		Os:         "linux",
		Filesystems: []*node.Filesystem{
			{Filesystem: "overlay", TypeName: "overlay", Total: 1023856, Free: 1023848, Available: 1023848, MountDir: "/"},
		},
	}

	err := printer.PrintNode(data, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testPrintPods(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	err := printer.PrintPods([]*pods.Pod{examplePod}, &buffer)
	assert.NoError(t, err, "Printing pods table should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func testPrintPod(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	data := &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      "foo",
			Namespace: "eliot",
		},
		Spec: &pods.PodSpec{
			Containers: []*containers.Container{
				{},
				{},
			},
		},
		Status: &pods.PodStatus{
			Hostname: "testing.local",
		},
	}

	err := printer.PrintPod(data, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
func testPrintConfig(t *testing.T, printer ResourcePrinter) {
	var buffer bytes.Buffer

	data := &config.Config{
		Endpoints: []config.Endpoint{
			{Name: "localhost", URL: "localhost:5000"},
		},
		Namespace: "default",
	}

	err := printer.PrintConfig(data, &buffer)
	assert.NoError(t, err, "Printing config should not return error")

	result := buffer.String()

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
