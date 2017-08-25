package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ernoaapa/layeryd/manifest"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/runtime"
	"github.com/ernoaapa/layeryd/state"
	"github.com/urfave/cli"
)

func getDeviceInfo() model.DeviceInfo {
	return model.DeviceInfo{}
}

func getRuntimeClient(clicontext *cli.Context) *runtime.ContainerdClient {
	return runtime.NewContainerdClient(
		context.Background(),
		clicontext.GlobalDuration("timeout"),
		clicontext.GlobalString("address"),
	)
}

func getManifestSource(clicontext *cli.Context) (manifest.Source, error) {
	file := clicontext.String("manifest-file")
	if file != "" {
		interval, err := time.ParseDuration(clicontext.String("manifest-update-interval"))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --manifest-update-interval 1s", clicontext.String("interval"))
		}
		return manifest.NewFileManifestSource(file, interval), nil
	}
	return nil, fmt.Errorf("You must define one manifest source for updates. E.g. --manifest-file path/to/file.yml")
}

func getStateReporter(clicontext *cli.Context, info model.DeviceInfo, client *runtime.ContainerdClient) (state.Reporter, error) {
	return state.NewConsoleStateReporter(
		info,
		client,
		clicontext.GlobalDuration("state-update-interval"),
	), nil
}
