package main

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/ernoaapa/layeryd/model"
	"github.com/ernoaapa/layeryd/source"
	"github.com/ernoaapa/layeryd/status"
	"github.com/urfave/cli"
)

func getDeviceInfo() model.DeviceInfo {
	return model.DeviceInfo{}
}

// appContext returns the context for a command. Should only be called once per
// command, near the start.
//
// This will ensure the namespace is picked up and set the timeout, if one is
// defined.
func appContext(clicontext *cli.Context) (context.Context, context.CancelFunc) {
	var (
		ctx       = context.Background()
		timeout   = clicontext.GlobalDuration("timeout")
		namespace = clicontext.GlobalString("namespace")
		cancel    context.CancelFunc
	)

	ctx = namespaces.WithNamespace(ctx, namespace)

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	return ctx, cancel
}

func getContainerdClient(clicontext *cli.Context) (*containerd.Client, error) {
	address := clicontext.GlobalString("address")
	namespace := clicontext.GlobalString("namespace")

	return containerd.New(address, containerd.WithDefaultNamespace(namespace))
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
