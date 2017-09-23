package printers

import (
	"io"

	"github.com/ernoaapa/can/pkg/model"
)

// ResourcePrinter is an interface that knows how to print runtime objects.
type ResourcePrinter interface {
	// Print receives a list of Pods, formats it and prints it to a writer.
	PrintPods([]*model.Pod, io.Writer) error
}
