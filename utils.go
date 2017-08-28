package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ernoaapa/layeryd/manifest"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/runtime"
	"github.com/ernoaapa/layeryd/state"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func getRuntimeClient(clicontext *cli.Context) *runtime.ContainerdClient {
	return runtime.NewContainerdClient(
		context.Background(),
		clicontext.GlobalDuration("timeout"),
		clicontext.GlobalString("address"),
	)
}

func getManifestSource(clicontext *cli.Context) (manifest.Source, error) {
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

func getStateReporter(clicontext *cli.Context, info *model.DeviceInfo, client *runtime.ContainerdClient) (state.Reporter, error) {
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
