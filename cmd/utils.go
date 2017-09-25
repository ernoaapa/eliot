package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ernoaapa/can/pkg/printers"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/can/pkg/api"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/controller"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/manifest"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// GetClient creates new cloud API client
func GetClient(clicontext *cli.Context) *api.Client {
	config := GetConfig(clicontext)
	return api.NewClient(
		config.GetCurrentContext().Namespace,
		config.GetCurrentEndpoint().URL,
	)
}

// GetConfig parse yaml config and return Config
func GetConfig(clicontext *cli.Context) *config.Config {
	configPath := clicontext.GlobalString("config")
	config, err := config.GetConfig(expandTilde(configPath))
	if err != nil {
		log.Fatalf("Error while reading configuration file [%s]: %s", configPath, err)
	}
	return config
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
		clicontext.GlobalString("address"),
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
