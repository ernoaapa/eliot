package printers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"sort"
	"strings"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/model"
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
	fmt.Fprintln(writer, "NAMESPACE\tNAME\tCONTAINERS\tSTATUS")

	for _, pod := range pods {
		fmt.Fprintf(writer, "%s\t%s\t%d\t%s", pod.Metadata.Name, pod.Metadata.Namespace, len(pod.Spec.Containers), getStatus(pod))
		fmt.Fprint(writer, "\n")
	}
	fmt.Fprint(writer, "\n")
	return nil
}

// getStatus constructs a string representation of all containers statuses
func getStatus(pod *pb.Pod) string {
	counts := map[string]int{}
	for _, status := range pod.Status.ContainerStatuses {
		if _, ok := counts[status.State]; !ok {
			counts[status.State] = 0
		}
		counts[status.State]++
	}

	keys := getKeys(counts)
	sort.Strings(keys)

	result := []string{}
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s(%d)", key, counts[key]))
	}
	return strings.Join(result, ",")
}

func getKeys(source map[string]int) (result []string) {
	for key := range source {
		result = append(result, key)
	}
	return result
}

// PrintDevicesTable writes stream of Devices in human readable table format to the writer
func (p *HumanReadablePrinter) PrintDevicesTable(devices <-chan model.DeviceInfo, writer io.Writer) error {
	fmt.Fprintln(writer, "HOSTNAME\tENDPOINT")

	go func(devices <-chan model.DeviceInfo) {
		for device := range devices {
			fmt.Fprintf(writer, "%s\t%s", device.Hostname, device.GetPrimaryEndpoint())
			fmt.Fprint(writer, "\n")
		}
	}(devices)
	return nil
}

// PrintPodDetails writes list of pods in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintPodDetails(pod *pb.Pod, writer io.Writer) error {
	t := template.New("pod-details")
	t, err := t.Parse(humanreadable.PodDetailsTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}
	data := map[string]interface{}{
		"Pod":    pod,
		"Status": getStatus(pod),
	}
	if err := t.Execute(writer, data); err != nil {
		return err
	}
	return nil
}

// PrintConfig writes list of pods in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintConfig(config *config.Config, writer io.Writer) error {
	t := template.New("config")
	t, err := t.Parse(humanreadable.ConfigTemplate)
	if err != nil {
		log.Fatalf("Invalid config template: %s", err)
	}

	if err := t.Execute(writer, config); err != nil {
		return err
	}

	fmt.Fprint(writer, "\n")
	return nil
}
