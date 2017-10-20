package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ernoaapa/can/pkg/discovery"
	"github.com/ernoaapa/can/pkg/printers"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/can/pkg/api"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/controller"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/fs"
	"github.com/ernoaapa/can/pkg/manifest"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	// GlobalFlags are flags what all commands have common
	GlobalFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
	}
)

// GlobalBefore is function what get executed before any commands executes
func GlobalBefore(context *cli.Context) error {
	if context.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

// GetClient creates new cloud API client
func GetClient(config *config.Provider) api.Client {
	return api.NewDirectClient(
		config.GetNamespace(),
		config.GetEndpoint(),
	)
}

// GetConfig parse yaml config and return the file representation
// In normal cases, you should use GetConfigProvider
func GetConfig(clicontext *cli.Context) *config.Config {
	configPath := clicontext.GlobalString("config")
	conf, err := config.GetConfig(expandTilde(configPath))
	if err != nil {
		log.Fatalf("Error while reading configuration file [%s]: %s", configPath, err)
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
		provider.OverrideEndpoint(clicontext.GlobalString("endpoint"))
	}

	if provider.GetEndpoint() == "" {
		log.Infoln("No endpoint url defined in configuration, try to discover from network automatically...")

		devices, err := discovery.Devices(2 * time.Second)
		if err != nil {
			log.Fatalf("Failed to auto-discover devices in network: %s", err)
		}

		if len(devices) == 1 {
			device := devices[0]
			log.Infof("Discovered device [%s] Will use endpoint [%s]", device.Hostname, device.GetPrimaryEndpoint())
			provider.OverrideEndpoint(devices[0].GetPrimaryEndpoint())
		} else if len(devices) > 1 {
			log.Fatalf("Discovered %d devices from network, get list of devices with command 'get devices'", len(devices))
		} else {
			log.Fatalf("Unable to discover device Automatically. You must give device endpoint. E.g. --endpoint=192.168.1.2")
		}
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
			log.Fatalf("Invalid --labels parameter [%s]. It must be comma separated key=value list. E.g. '--labels foo=bar,one=two'", param)
		}
	}
	return labels
}

// GetRuntimeClient initialises new runtime client from CLI parameters
func GetRuntimeClient(clicontext *cli.Context) runtime.Client {
	return runtime.NewContainerdClient(
		context.Background(),
		clicontext.GlobalDuration("timeout"),
		clicontext.GlobalString("containerd"),
	)
}

// GetManifestSource initialises new manifest source based on CLI parameters
func GetManifestSource(clicontext *cli.Context, resolver *device.Resolver, out chan<- []model.Pod) (manifest.Source, error) {
	if !clicontext.IsSet("manifest") {
		return nil, fmt.Errorf("You must define --manifest parameter")
	}
	manifestParam := clicontext.String("manifest")

	interval := clicontext.Duration("manifest-update-interval")

	if fileExists(manifestParam) {
		return manifest.NewFileManifestSource(manifestParam, interval, resolver, out), nil
	}

	manifestURL, err := url.Parse(manifestParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --manifest parameter [%s]", manifestParam)
	}

	switch scheme := manifestURL.Scheme; scheme {
	case "file":
		return manifest.NewFileManifestSource(manifestURL.Path, interval, resolver, out), nil
	case "http", "https":
		return manifest.NewURLManifestSource(manifestParam, interval, resolver, out), nil
	}
	return nil, fmt.Errorf("You must define manifest source. E.g. --manifest path/to/file.yml")
}

// GetStateReporter initialises new state reporter based on CLI parameters
func GetStateReporter(clicontext *cli.Context, resolver *device.Resolver, in <-chan []model.Pod) (state.Reporter, error) {
	if !clicontext.IsSet("report") {
		return nil, fmt.Errorf("You must define --report parameter. E.g. 'console' or some url")
	}
	reportParam := clicontext.String("report")
	if reportParam == "console" {
		return state.NewConsoleStateReporter(
			resolver,
			in,
		), nil
	}
	_, err := url.Parse(reportParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --report parameter [%s]", reportParam)
	}

	return state.NewURLStateReporter(
		resolver,
		in,
		reportParam,
	), nil
}

// GetController creates new Controller
func GetController(clicontext *cli.Context, in <-chan []model.Pod, out chan<- []model.Pod) *controller.Controller {
	client := GetRuntimeClient(clicontext)
	interval := clicontext.Duration("update-interval")
	return controller.New(client, interval, in, out)
}

// GetPrinter returns printer for formating resources output
func GetPrinter(clicontext *cli.Context) printers.ResourcePrinter {
	return printers.NewHumanReadablePrinter()
}

// GetMounts parses a --mount string flags
func GetMounts(clicontext *cli.Context) (result []*containers.Mount) {
	for _, flag := range clicontext.StringSlice("mount") {
		mount, err := parseMountFlag(flag)
		if err != nil {
			log.Fatalf("Failed to parse --mount flag: %s", err)
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

// GetBinds parses a --bind string flags
func GetBinds(clicontext *cli.Context, extra ...string) (result []*containers.Mount) {
	binds := clicontext.StringSlice("bind")
	binds = append(binds, extra...)
	for _, flag := range binds {
		bind, err := ParseBindFlag(flag)
		if err != nil {
			log.Fatalf("Failed to parse --bind flag: %s", err)
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func expandTilde(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path[:2] == "~/" {
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

var (
	defaultRegistry = "docker.io"
	defaultUsername = "library"
	defaultTag      = "latest"
)

// ExpandToFQIN converts partial image name to "Fully Qualified Image Name"
// E.g. eaapa/hello-world -> docker.io/eaapa/hello-world:latest
func ExpandToFQIN(source string) string {
	registry := defaultRegistry
	username := defaultUsername
	tag := defaultTag
	image := source

	parts := strings.SplitN(source, "/", 3)
	if len(parts) == 3 {
		registry = parts[0]
		username = parts[1]
		image = parts[2]
	} else if len(parts) == 2 {
		username = parts[0]
		image = parts[1]
	}

	imageParts := strings.SplitN(image, ":", 2)
	if len(imageParts) == 2 {
		image = imageParts[0]
		tag = imageParts[1]
	}

	return fmt.Sprintf("%s/%s/%s:%s", registry, username, image, tag)
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
				log.WithError(err).Errorf("forward signal %s", s)
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

	log.Fatal("Failed to resolve current directory")
	return ""
}
