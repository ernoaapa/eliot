package printers

import (
	"html/template"
	"io"
	"log"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
)

// HumanReadablePrinter is an implementation of ResourcePrinter which prints
// resources in human readable format (tables etc.).
type HumanReadablePrinter struct {
}

// NewHumanReadablePrinter creates new HumanReadablePrinter
func NewHumanReadablePrinter() *HumanReadablePrinter {
	return &HumanReadablePrinter{}
}

const podTemplate = `NAMESPACE	NAME	CONTAINERS{{range .}}
{{.Metadata.namespace}}	{{.Metadata.name}}	{{len .Spec.Containers}}
{{else}}
No pods
{{end}}
`

// PrintPods writes list of Pods in human readable format to the writer
func (p *HumanReadablePrinter) PrintPods(data []*pb.Pod, writer io.Writer) error {
	t := template.New("pods")
	t, err := t.Parse(podTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}

	if err := t.Execute(writer, data); err != nil {
		return err
	}
	return nil
}
