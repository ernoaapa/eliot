package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ernoaapa/layeryd/manifest"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/runtime"
	"github.com/ernoaapa/layeryd/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

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

func GetRuntimeClient(clicontext *cli.Context) *runtime.ContainerdClient {
	return runtime.NewContainerdClient(
		context.Background(),
		clicontext.GlobalDuration("timeout"),
		clicontext.GlobalString("address"),
	)
}

func GetManifestSource(clicontext *cli.Context) (manifest.Source, error) {
	if !clicontext.IsSet("manifest") {
		return nil, fmt.Errorf("You must define --manifest parameter")
	}
	manifestParam := clicontext.String("manifest")

	interval, err := time.ParseDuration(clicontext.String("manifest-update-interval"))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --manifest-update-interval 1s", clicontext.String("interval"))
	}

	if fileExists(manifestParam) {
		return manifest.NewFileManifestSource(manifestParam, interval), nil
	}

	manifestURL, err := url.Parse(manifestParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --manifest parameter [%s]", manifestParam)
	}

	switch scheme := manifestURL.Scheme; scheme {
	case "file":
		return manifest.NewFileManifestSource(manifestURL.Path, interval), nil
	case "http", "https":
		return manifest.NewURLManifestSource(manifestParam, interval), nil
	}
	return nil, fmt.Errorf("You must define manifest source. E.g. --manifest path/to/file.yml")
}

func GetStateReporter(clicontext *cli.Context, info *model.DeviceInfo, client *runtime.ContainerdClient) (state.Reporter, error) {
	return state.NewConsoleStateReporter(
		info,
		client,
		clicontext.GlobalDuration("state-update-interval"),
	), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
