package printers

import (
	"fmt"
	"html/template"
	"io"
	"log"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/printers/humanreadable"
)

// HumanReadablePrinter is an implementation of ResourcePrinter which prints
// resources in human readable format (tables etc.).
type HumanReadablePrinter struct {
}

// NewHumanReadablePrinter creates new HumanReadablePrinter
func NewHumanReadablePrinter() *HumanReadablePrinter {
	return &HumanReadablePrinter{}
}

// PrintPodsTable writes list of Pods in human readable table format to the writer
func (p *HumanReadablePrinter) PrintPodsTable(pods []*pb.Pod, writer io.Writer) error {
	fmt.Fprintln(writer, humanreadable.PodsTableHeader)
	t := template.New("pods-table")
	t, err := t.Parse(humanreadable.PodsTableRowTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}

	for _, pod := range pods {
		if err := t.Execute(writer, pod); err != nil {
			return err
		}
	}
	// TODO: For some reason, don't output without printing something to the writer
	// Find out how to flush the writer
	fmt.Fprintln(writer, "\n")
	return nil
}

// PrintPodDetails writes list of pods in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintPodDetails(pod *pb.Pod, writer io.Writer) error {
	t := template.New("pod-details")
	t, err := t.Parse(humanreadable.PodDetailsTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}

	if err := t.Execute(writer, pod); err != nil {
		return err
	}
	return nil
}
