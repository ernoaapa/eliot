package printers

import (
	"io"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
)

// ResourcePrinter is an interface that knows how to print runtime objects.
type ResourcePrinter interface {
	// Print receives a list of Pods, formats it as table and prints it to a writer.
	PrintPodsTable([]*pb.Pod, io.Writer) error
	// Print receives a list of Pods, formats it to detailed description and prints it to the writer.
	PrintPodDetails(*pb.Pod, io.Writer) error
}
