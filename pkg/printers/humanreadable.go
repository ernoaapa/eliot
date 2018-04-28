package printers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/printers/humanreadable"
	"github.com/ernoaapa/eliot/pkg/utils"
	"github.com/pkg/errors"

	"github.com/c2h5oh/datasize"
	"github.com/hako/durafmt"
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

// PrintNodes writes list of Nodes in human readable table format to the writer
func (p *HumanReadablePrinter) PrintNodes(nodes []*node.Info, writer io.Writer) error {
	if len(nodes) == 0 {
		fmt.Fprintf(writer, "\n\t(No nodes)\n\n")
		return nil
	}
	fmt.Fprintln(writer, "\nHOSTNAME\tENDPOINT\tVERSION")

	for _, node := range nodes {
		endpoint := fmt.Sprintf("%s:%d", utils.GetFirst(node.Addresses, ""), node.GrpcPort)
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", node.Hostname, endpoint, node.Version)
		if err != nil {
			return errors.Wrapf(err, "Error while writing node row")
		}
	}

	return nil
}

// PrintNode writes a node in human readable detailed format to the writer
func (p *HumanReadablePrinter) PrintNode(info *node.Info, writer io.Writer) error {
	t := template.New("node-details").Funcs(template.FuncMap{
		"FormatPercent": formatPercent,
		"FormatUptime":  formatUptime,
		"Subtract": func(a, b uint64) uint64 {
			return a - b
		},
		"FormatBytes": func(v uint64) string {
			return datasize.ByteSize(v).HumanReadable()
		},
	})
	t, err := t.Parse(humanreadable.NodeDetailsTemplate)
	if err != nil {
		log.Fatalf("Invalid pod template: %s", err)
	}
	return t.Execute(writer, info)
}

func formatPercent(total, free, available uint64) string {
	percent := 0.0
	bUsed := (total - free) / 1024
	bAvail := available / 1024
	utotal := bUsed + bAvail
	used := bUsed

	if utotal != 0 {
		u100 := used * 100
		pct := u100 / utotal
		if u100%utotal != 0 {
			pct++
		}
		percent = (float64(pct) / float64(100)) * 100.0
	}

	return strconv.FormatFloat(percent, 'f', -1, 64) + "%"
}

func formatUptime(uptime uint64) string {
	var duration = time.Duration(uptime * 1000 * 1000 * 1000)
	if duration < 0 {
		// the duration went over maximum int64, fallback to just display the seconds
		return fmt.Sprintf("%d seconds", uptime)
	}
	return durafmt.Parse(duration).String()
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
