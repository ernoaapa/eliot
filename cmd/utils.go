package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ernoaapa/can/pkg/client"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/manifest"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func GetClient(clicontext *cli.Context) *client.Client {
	configPath := clicontext.GlobalString("config")
	config, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Error while reading configuration file [%s]: %s", configPath, err)
	}
	return client.NewClient(
		config.GetCurrentEndpoint().URL,
		config.GetCurrentUser().Token,
	)
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

	interval, err := time.ParseDuration(clicontext.String("manifest-update-interval"))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --manifest-update-interval 1s", clicontext.String("interval"))
	}

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
func GetStateReporter(clicontext *cli.Context, resolver *device.Resolver, client runtime.Client) (state.Reporter, error) {
	reportParam := clicontext.String("report")
	if reportParam == "console" {
		return state.NewConsoleStateReporter(
			resolver,
			client,
			clicontext.GlobalDuration("state-update-interval"),
		), nil
	}
	_, err := url.Parse(reportParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --report parameter [%s]", reportParam)
	}

	return state.NewURLStateReporter(
		resolver,
		client,
		clicontext.GlobalDuration("state-update-interval"),
		reportParam,
	), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
