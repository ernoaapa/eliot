package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ernoaapa/eliot/pkg/cmd"
	ui "github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/discovery"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/ernoaapa/eliot/pkg/sync"
	"github.com/ernoaapa/eliot/pkg/utils"

	"github.com/sirupsen/logrus"

	"github.com/ernoaapa/eliot/pkg/api"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/fs"
	"github.com/ernoaapa/eliot/pkg/runtime"
	"github.com/urfave/cli"
)

const (
	outputHuman = "human"
	outputYaml  = "yaml"
)

var (
	// GlobalFlags are flags what all commands have common
	GlobalFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "Don't print any progress output",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: fmt.Sprintf("Output format. One of: %s", []string{outputHuman, outputYaml}),
			Value: "human",
		},
	}
)

// GlobalBefore is function what get executed before any commands executes
func GlobalBefore(context *cli.Context) error {
	debug := context.GlobalBool("debug")
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cmd.IsPipingOut() || context.GlobalBool("quiet") || context.GlobalString("output") != outputHuman {
		ui.SetOutput(ui.NewHidden())
	} else if debug {
		ui.SetOutput(ui.NewDebug())
	} else {
		ui.SetOutput(ui.NewTerminal())
	}

	return nil
}

// GetClient creates new cloud API client
func GetClient(config *config.Provider) *api.Client {
	uiline := ui.NewLine()

	endpoints := config.GetEndpoints()
	switch len(endpoints) {
	case 0:
		uiline.Fatal("No node to connect. You must give node endpoint. E.g. --endpoint=192.168.1.2")
		return nil
	case 1:
		uiline.Loadingf("Connecting to %s (%s)", endpoints[0].Name, endpoints[0].URL)
		client := api.NewClient(config.GetNamespace(), endpoints[0])
		info, err := client.GetInfo()
		if err != nil {
			logrus.Debugf("Connection failure: %s", err)
			uiline.Fatalf("Failed connect to %s (%s)", endpoints[0].Name, endpoints[0].URL)
		}
		uiline.Donef("Connected to %s (%s)", info.Hostname, endpoints[0].URL)
		return client
	default:
		uiline.Fatalf("%d node found. You must give target node. E.g. --endpoint=192.168.1.2", len(endpoints))
		return nil
	}
}

// GetConfig parse yaml config and return the file representation
// In normal cases, you should use GetConfigProvider
func GetConfig(clicontext *cli.Context) *config.Config {
	configPath := clicontext.GlobalString("config")
	conf, err := config.GetConfig(expandTilde(configPath))
	if err != nil {
		ui.NewLine().Fatalf("Error while reading configuration file [%s]: %s", configPath, err)
	}
	return conf
}

// GetConfigProvider return config.Provider to access the current configuration
func GetConfigProvider(clicontext *cli.Context) *config.Provider {

	provider := config.NewProvider(GetConfig(clicontext))

	if clicontext.GlobalIsSet("namespace") && clicontext.GlobalString("namespace") != "" {
		provider.OverrideNamespace(clicontext.GlobalString("namespace"))
	}

	if clicontext.GlobalIsSet("endpoint") && clicontext.GlobalString("endpoint") != "" {
		provider.OverrideEndpoints([]config.Endpoint{{
			Name: clicontext.GlobalString("endpoint"),
			URL:  clicontext.GlobalString("endpoint"),
		}})
	}

	if len(provider.GetEndpoints()) == 0 {
		uiline := ui.NewLine().Loading("Discover from network automatically...")
		node, err := discovery.Nodes(2 * time.Second)
		if err != nil {
			uiline.Errorf("Failed to auto-discover node in network: %s", err)
		} else {
			if len(node) == 0 {
				uiline.Warn("No node discovered from network")
			} else {
				uiline.Donef("Discovered %d node(s) from network", len(node))
			}
		}

		endpoints := []config.Endpoint{}
		for _, node := range node {
			if len(node.Addresses) > 0 {
				endpoints = append(endpoints, config.Endpoint{
					Name: node.Hostname,
					URL:  fmt.Sprintf("%s:%d", utils.GetFirst(node.Addresses, ""), node.GrpcPort),
				})
			}
		}
		provider.OverrideEndpoints(endpoints)
	}

	if clicontext.GlobalIsSet("node") && clicontext.GlobalString("node") != "" {
		nodeName := clicontext.GlobalString("node")
		endpoint, found := provider.GetEndpointByName(nodeName)
		if !found {
			ui.NewLine().Errorf("Failed to find node with name %s", nodeName)
		}
		provider.OverrideEndpoints([]config.Endpoint{endpoint})
	}

	return provider
}

// UpdateConfig writes config to the config file in yaml format
func UpdateConfig(clicontext *cli.Context, updated *config.Config) error {
	configPath := expandTilde(clicontext.GlobalString("config"))

	return config.WriteConfig(configPath, updated)
}

// GetLabels return --labels CLI parameter value as string map
func GetLabels(clicontext *cli.Context) map[string]string {
	if !clicontext.IsSet("labels") {
		return map[string]string{}
	}

	param := clicontext.String("labels")
	values := strings.Split(param, ",")

	labels := map[string]string{}
	for _, value := range values {
		pair := strings.Split(value, "=")
		if len(pair) == 2 {
			labels[pair[0]] = pair[1]
		} else {
			ui.NewLine().Fatalf("Invalid --labels parameter [%s]. It must be comma separated key=value list. E.g. '--labels foo=bar,one=two'", param)
		}
	}
	return labels
}

// GetRuntimeClient initialises new runtime client from CLI parameters
func GetRuntimeClient(clicontext *cli.Context, hostname string) runtime.Client {
	return runtime.NewContainerdClient(
		context.Background(),
		clicontext.GlobalDuration("timeout"),
		clicontext.String("containerd-snapshotter"),
		clicontext.GlobalString("containerd"),
		hostname,
	)
}

// GetPrinter returns printer for formating resources output
func GetPrinter(clicontext *cli.Context) printers.ResourcePrinter {
	switch output := clicontext.GlobalString("output"); output {
	case outputHuman:
		return printers.NewHumanReadablePrinter()
	case outputYaml:
		return printers.NewYamlPrinter()
	default:
		logrus.Fatalf("Unknown output format: %s", output)
		return nil
	}
}

// MustParseMounts parses a --mount string flags
func MustParseMounts(mounts []string) (result []*containers.Mount) {
	for _, str := range mounts {
		mount, err := parseMountFlag(str)
		if err != nil {
			ui.NewLine().Fatalf("Failed to parse --mount flag [%s]: %s", str, err)
		}
		result = append(result, mount)
	}
	return result
}

// parseMountFlag parses a mount string in the form "type=foo,source=/path,destination=/target,options=rbind:rw"
func parseMountFlag(m string) (*containers.Mount, error) {
	mount := &containers.Mount{}
	r := csv.NewReader(strings.NewReader(m))

	fields, err := r.Read()
	if err != nil {
		return mount, err
	}

	for _, field := range fields {
		v := strings.Split(field, "=")
		if len(v) != 2 {
			return mount, fmt.Errorf("invalid mount specification: expected key=val")
		}

		key := v[0]
		val := v[1]
		switch key {
		case "type":
			mount.Type = val
		case "source", "src":
			mount.Source = val
		case "destination", "dst":
			mount.Destination = val
		case "options":
			mount.Options = strings.Split(val, ":")
		default:
			return mount, fmt.Errorf("mount option %q not supported", key)
		}
	}

	return mount, nil
}

// MustParseBinds parses a --bind string flags
func MustParseBinds(binds []string) (result []*containers.Mount) {
	for _, flag := range binds {
		bind, err := ParseBindFlag(flag)
		if err != nil {
			ui.NewLine().Fatalf("Failed to parse --bind flag: %s", err)
		}
		result = append(result, bind)
	}
	return result
}

// MustParseBindFlag is like ParseBindFlag but panics if syntax is invalid
func MustParseBindFlag(b string) *containers.Mount {
	m, err := ParseBindFlag(b)
	if err != nil {
		panic("Invalid mount format: " + b + ". Error: " + err.Error())
	}
	return m
}

// ParseBindFlag parses a mount string in the form "/var:/var:rshared"
func ParseBindFlag(b string) (*containers.Mount, error) {
	parts := strings.Split(b, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("Cannot parse bind, missing ':': %s", b)
	}
	if len(parts) > 3 {
		return nil, fmt.Errorf("Cannot parse bind, too many ':': %s", b)
	}
	src := parts[0]
	dest := parts[1]
	opts := []string{"rw", "rbind", "rprivate"}
	if len(parts) == 3 {
		opts = append(strings.Split(parts[2], ","), "rbind")
	}
	return &containers.Mount{
		Type:        "bind",
		Destination: dest,
		Source:      src,
		Options:     opts,
	}, nil
}

// MustParseSyncs parses a sync string in the form "~/local/dir:/data"
func MustParseSyncs(syncs []string) (result []sync.Sync) {
	for _, value := range syncs {
		sync, err := sync.Parse(value)
		if err != nil {
			logrus.Fatalf("Error reading sync argument: %s", err)
		}
		result = append(result, sync)
	}

	return result
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func expandTilde(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if len(path) >= 2 && path[:2] == "~/" {
		return filepath.Join(dir, path[2:])
	}
	return path
}

// FilterByPodName return new list of Pods which name matches with given podName
func FilterByPodName(source []*pods.Pod, podName string) []*pods.Pod {
	for _, pod := range source {
		if pod.Metadata.Name == podName {
			return []*pods.Pod{
				pod,
			}
		}
	}
	return []*pods.Pod{}
}

// ForwardAllSignals will listen all kill signals and pass it to the handler
func ForwardAllSignals(handler func(syscall.Signal) error) chan os.Signal {
	sigc := make(chan os.Signal, 128)
	signal.Notify(sigc)
	go func() {
		for s := range sigc {
			signal := s.(syscall.Signal)
			// Doesn't make sense to forward "child process terminates" because it's about this CLI child process
			if signal == syscall.SIGCHLD {
				continue
			}

			if err := handler(signal); err != nil {
				logrus.WithError(err).Errorf("forward signal %s", s)
			}
		}
	}()
	return sigc
}

// GetCurrentDirectory resolves current directory where the command were executed
// Tries different options until find one or fails
func GetCurrentDirectory() string {
	for _, path := range []string{".", os.Args[0], os.Getenv("PWD")} {
		dir, err := filepath.Abs(filepath.Dir(path))
		if err == nil && fs.DirExist(path) {
			return dir
		}
	}

	ui.NewLine().Fatal("Failed to resolve current directory")
	return ""
}

// First return first non empty "" string or empty ""
func First(values ...string) string {
	for _, str := range values {
		if str != "" {
			return str
		}
	}
	return ""
}

// DropDoubleDash search for double dash (--) and if found
// return arguments after it, otherwise return all arguments
func DropDoubleDash(args []string) []string {
	for index, arg := range args {
		if arg == "--" {
			return args[index+1:]
		}
	}
	return args
}

// ResolveContainerID resolves ContainerID from list of containers.
// If multiple containers, you must define containerName, otherwise it's optional.
func ResolveContainerID(containers []*containers.ContainerStatus, containerName string) (string, error) {
	containerCount := len(containers)
	if containerCount == 0 {
		return "", fmt.Errorf("Pod don't have any containers")
	} else if containerCount == 1 {
		return containers[0].ContainerID, nil
	} else {
		if containerName == "" {
			return "", fmt.Errorf("Pod contains %d containers, you must define container name", containerCount)
		}

		for _, status := range containers {
			if status.Name == containerName {
				return status.ContainerID, nil
			}
		}

		return "", fmt.Errorf("Pod contains %d containers, you must define container name", containerCount)
	}
}

// FindRunningContainerID search from Pod definition a containerID by container name
func FindRunningContainerID(pod *pods.Pod, name string) (string, error) {
	if pod.Status != nil && len(pod.Status.ContainerStatuses) > 0 {
		for _, status := range pod.Status.ContainerStatuses {
			if status.Name == name {
				return status.ContainerID, nil
			}
		}
	}

	return "", fmt.Errorf("Cannot find ContainerID with name %s", name)
}

// StopCatch will close the given channel when receives Stop signal (^C)
func StopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}
