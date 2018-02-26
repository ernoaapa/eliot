package printers

import (
	"io"

	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
)

// ResourcePrinter is an interface that knows how to print runtime objects.
type ResourcePrinter interface {
	PrintPods([]*pods.Pod, io.Writer) error
	PrintNodes([]*node.Info, io.Writer) error
	PrintNode(*node.Info, io.Writer) error
	PrintPod(*pods.Pod, io.Writer) error
	PrintConfig(*config.Config, io.Writer) error
}
