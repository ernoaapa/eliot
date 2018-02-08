package printers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"sort"
	"strings"

	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/printers/humanreadable"
	"github.com/pkg/errors"
)

// HumanReadablePrinter is an implementation of ResourcePrinter which prints
// resources in human readable format (tables etc.).
type HumanReadablePrinter struct {
}

// NewHumanReadablePrinter creates new HumanReadablePrinter
func NewHumanReadablePrinter() *HumanReadablePrinter {
	return &HumanReadablePrinter{}
}

// PrintPods writes list of Pods in human readable table format to the writer
func (p *HumanReadablePrinter) PrintPods(pods []*pods.Pod, writer io.Writer) error {
	if len(pods) == 0 {
		fmt.Fprintf(writer, "\n\t(No pods)\n\n")
		return nil
	}

	fmt.Fprintln(writer, "\nNAMESPACE\tNAME\tCONTAINERS\tSTATUS")

	for _, pod := range pods {
		_, err := fmt.Fprintf(writer, "%s\t%s\t%d\t%s\n", pod.Metadata.Namespace, pod.Metadata.Name, len(pod.Spec.Containers), getStatus(pod))
		if err != nil {
			return errors.Wrapf(err, "Error while writing pod row")
		}
	}

	return nil
}

// getStatus constructs a string representation of all containers statuses
func getStatus(pod *pods.Pod) string {
	counts := map[string]int{}

	statuses := []*containers.ContainerStatus{}
	if pod.Status != nil {
		statuses = pod.Status.ContainerStatuses
	}
	for _, status := range statuses {
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

// PrintDevices writes list of Devices in human readable table format to the writer
func (p *HumanReadablePrinter) PrintDevices(devices []model.DeviceInfo, writer io.Writer) error {
	if len(devices) == 0 {
		fmt.Fprintf(writer, "\n\t(No devices)\n\n")
		return nil
	}
	fmt.Fprintln(writer, "\nHOSTNAME\tENDPOINT\tVERSION")

	for _, device := range devices {
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", device.Hostname, device.GetPrimaryEndpoint(), device.Version)
		if err != nil {
			return errors.Wrapf(err, "Error while writing device row")
		}
	}

	return nil
}

// PrintDevice writes a device in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintDevice(info *device.Info, writer io.Writer) error {
	t := template.New("device-details")
	t, err := t.Parse(humanreadable.DeviceDetailsTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}
	return t.Execute(writer, info)
}

// PrintPod writes a pod in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintPod(pod *pods.Pod, writer io.Writer) error {
	t := template.New("pod-details").Funcs(template.FuncMap{
		"GetStatus": func(pod pods.Pod, name string) *containers.ContainerStatus {
			if pod.Status == nil {
				return nil
			}
			for _, status := range pod.Status.ContainerStatuses {
				if status.Name == name {
					return status
				}
			}
			return nil
		},
		"StringsJoin": strings.Join,
	})
	t, err := t.Parse(humanreadable.PodDetailsTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}
	data := map[string]interface{}{
		"Pod":    pod,
		"Status": getStatus(pod),
	}
	return t.Execute(writer, data)
}

// PrintConfig writes list of pods in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintConfig(config *config.Config, writer io.Writer) error {
	t := template.New("config")
	t, err := t.Parse(humanreadable.ConfigTemplate)
	if err != nil {
		log.Fatalf("Invalid config template: %s", err)
	}

	return t.Execute(writer, config)
}
