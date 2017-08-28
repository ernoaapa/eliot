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
	manifestParam := clicontext.String("manifest")
	if manifestParam == "" {
		return nil, fmt.Errorf("You must define --manifest parameter")
	}

	if fileExists(manifestParam) {
		return getFileManifestSource(clicontext, manifestParam)
	}

	manifestURL, err := url.Parse(manifestParam)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while parsing --manifest parameter [%s]", manifestParam)
	}

	switch scheme := manifestURL.Scheme; scheme {
	case "file":
		return getFileManifestSource(clicontext, manifestURL.Path)
	}
	return nil, fmt.Errorf("You must define manifest source. E.g. --manifest path/to/file.yml")
}

func getFileManifestSource(clicontext *cli.Context, path string) (manifest.Source, error) {
	interval, err := time.ParseDuration(clicontext.String("manifest-update-interval"))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --manifest-update-interval 1s", clicontext.String("interval"))
	}
	return manifest.NewFileManifestSource(path, interval), nil
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
