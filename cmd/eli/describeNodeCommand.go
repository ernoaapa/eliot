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

var describeNodeCommand = cli.Command{
	Name:    "node",
	Aliases: []string{"nodes", "device"}, // device is deprecated command
	Usage:   "Return details of node",
	UsageText: `eli describe node [options] [NAME | IP]
	
	# Describe node
	eli describe node my-node

	# Describe all pods
	eli describe nodes
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
				return errors.Wrap(err, "Failed to fetch node info")
			}
			if err := printer.PrintNode(info, writer); err != nil {
				return err
			}
		}
		return nil
	},
}
