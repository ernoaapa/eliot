package printers

import (
	"bytes"
	"testing"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var buffer bytes.Buffer
	printer := NewHumanReadablePrinter()

	data := []*pb.Pod{
		&pb.Pod{
			Metadata: map[string]string{
				"name":      "foo",
				"namespace": "cand",
			},
			Spec: &pb.PodSpec{
				Containers: []*pb.Container{
					&pb.Container{},
					&pb.Container{},
				},
			},
		},
	}

	printer.PrintPods(data, &buffer)

	result := buffer.String()

	log.Debugln(result)

	assert.True(t, len(result) > 0, "Should write something to the writer")
}
