package printers

import (
	"bytes"
	"testing"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPrintTable(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := []*pb.Pod{
		&pb.Pod{
			Metadata: &pb.ResourceMetadata{
				Name:      "foo",
				Namespace: "cand",
			},
			Spec: &pb.PodSpec{
				Containers: []*pb.Container{
					&pb.Container{},
					&pb.Container{},
				},
			},
		},
	}

	err := printer.PrintPodsTable(data, &buffer)
	assert.NoError(t, err, "Printing pods table should not return error")

	result := buffer.String()

	log.Debugln(result)

	assert.True(t, len(result) > 0, "Should write something to the writer")
}

func TestPrintDetails(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := &pb.Pod{
		Metadata: &pb.ResourceMetadata{
			Name:      "foo",
			Namespace: "cand",
		},
		Spec: &pb.PodSpec{
			Containers: []*pb.Container{
				&pb.Container{},
				&pb.Container{},
			},
		},
	}

	err := printer.PrintPodDetails(data, &buffer)
	assert.NoError(t, err, "Printing pod details should not return error")

	result := buffer.String()

	log.Debugln(result)

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
