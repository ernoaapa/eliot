package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var describeDeviceCommand = cli.Command{
	Name:    "device",
	Aliases: []string{"devices"},
	Usage:   "Return details of device",
	UsageText: `eli describe device [options] [NAME | IP]
	
	# Describe device
	eli describe device my-device

	# Describe all pods
	eli describe devices
`,
	Action: func(clicontext *cli.Context) error {
		cfg := cmd.GetConfigProvider(clicontext)

		podName := clicontext.Args().First()

		endpoints := cfg.GetEndpoints()

		if podName != "" {
			endpoint, ok := cfg.GetEndpointByName(podName)
			if !ok {
				return fmt.Errorf("No endpoint found with name %s", podName)
			}
			endpoints = []config.Endpoint{endpoint}
		}

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()
		printer := cmd.GetPrinter(clicontext)
		for _, endpoint := range endpoints {
			info, err := api.NewClient(cfg.GetNamespace(), endpoint).GetInfo()
			if err != nil {
				return errors.Wrap(err, "Failed to fetch device info")
			}
			if err := printer.PrintDevice(info, writer); err != nil {
				return err
			}
		}
		return nil
	},
}
