package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ernoaapa/layeryd/runtime"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/source"
	"github.com/ernoaapa/layeryd/status"
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
		clicontext.GlobalString("namespace"),
	)
}

func getSource(clicontext *cli.Context) (source.Source, error) {
	file := clicontext.String("file")
	if file != "" {
		interval, err := time.ParseDuration(clicontext.String("interval"))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse update interval [%s]. Example --interval 1s", clicontext.String("interval"))
		}
		return source.NewFileSource(file, interval), nil
	}
	return nil, fmt.Errorf("You must define one source for updates. E.g. --file path/to/file.yml")
}

func getReporter(clicontext *cli.Context) (status.Reporter, error) {
	return status.NewConsoleReporter(), nil
}
