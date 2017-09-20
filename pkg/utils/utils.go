package utils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/manifest"
	"github.com/ernoaapa/can/pkg/runtime"
	"github.com/ernoaapa/can/pkg/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

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
func GetManifestSource(clicontext *cli.Context, resolver *device.Resolver) (manifest.Source, error) {
	if !clicontext.IsSet("manifest") {
		return nil, fmt.Errorf("You must define --manifest parameter")
	}
	manifestParam := clicontext.String("manifest")

	interval, err := time.ParseDuration(clicontext.String("manifest-update-interval"))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --manifest-update-interval 1s", clicontext.String("interval"))
	}

	if fileExists(manifestParam) {
		return manifest.NewFileManifestSource(manifestParam, interval, resolver), nil
	}

	manifestURL, err := url.Parse(manifestParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --manifest parameter [%s]", manifestParam)
	}

	switch scheme := manifestURL.Scheme; scheme {
	case "file":
		return manifest.NewFileManifestSource(manifestURL.Path, interval, resolver), nil
	case "http", "https":
		return manifest.NewURLManifestSource(manifestParam, interval, resolver), nil
	}
	return nil, fmt.Errorf("You must define manifest source. E.g. --manifest path/to/file.yml")
}

// GetStateReporter initialises new state reporter based on CLI parameters
func GetStateReporter(clicontext *cli.Context, resolver *device.Resolver, client runtime.Client) (state.Reporter, error) {
	// return state.NewConsoleStateReporter(
	// 	resolver,
	// 	client,
	// 	clicontext.GlobalDuration("state-update-interval"),
	// ), nil
	return state.NewURLStateReporter(
		resolver,
		client,
		clicontext.GlobalDuration("state-update-interval"),
		"http://localhost:3000/api/devices/localtest/pods",
	), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
