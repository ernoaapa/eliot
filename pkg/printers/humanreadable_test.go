package printers

import (
	"bytes"
	"testing"

	"github.com/ernoaapa/can/pkg/api/core"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPrintTable(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := []*pods.Pod{
		&pods.Pod{
			Metadata: &core.ResourceMetadata{
				Name:      "foo",
				Namespace: "cand",
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

func TestPrintDetails(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      "foo",
			Namespace: "cand",
		},
		Spec: &pods.PodSpec{
			Containers: []*containers.Container{
				&containers.Container{},
				&containers.Container{},
			},
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
